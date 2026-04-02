package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

/*
go build -o loadtest ./

	./loadtest \
	  -addr 10.1.0.119:19192 \
	  -agents 10 \
	  -duration 5m \
	  -rampup 30s \
	  -rate-ms 500 \
	  -containers 5

	  ./loadtest -addr 10.1.0.119:19192 -agents 1000 -duration 5m -rampup 30s -rate-ms 500 -containers 5
*/
func main() {
	addr := flag.String("addr", "10.1.0.119:19192", "gRPC server address")
	agentCount := flag.Int("agents", 1000, "number of load-test agents (IDs: 2 ~ agents+1)")
	duration := flag.Duration("duration", 5*time.Minute, "test duration (0 = run until SIGINT)")
	rampup := flag.Duration("rampup", 30*time.Second, "ramp-up period for agent start") // 30s / 1000 = 30ms 간격으로 에이전트를 순차 기동해서 커넥션을 분산
	rateMs := flag.Int("rate-ms", 500, "message generation interval per agent in ms")
	containers := flag.Int("containers", 5, "fake container count per agent")
	flag.Parse()

	metrics := &Metrics{}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	fmt.Printf("=== gRPC Load Test ===\n")
	fmt.Printf("addr=%s  agents=%d  agentId=[2..%d]\n", *addr, *agentCount, *agentCount+1)
	fmt.Printf("duration=%v  rampup=%v  rate=%dms  containers=%d\n\n",
		*duration, *rampup, *rateMs, *containers)

	var (
		activeAgents atomic.Int32
		agentWg      sync.WaitGroup
	)

	// Spread agent startup evenly across the ramp-up window.
	rampInterval := time.Duration(0)
	if *agentCount > 1 && *rampup > 0 {
		rampInterval = *rampup / time.Duration(*agentCount)
	}

	startTime := time.Now()
	started := 0

	for i := 0; i < *agentCount; i++ {
		select {
		case <-ctx.Done():
			fmt.Printf("Interrupted during ramp-up (%d/%d agents started)\n", started, *agentCount)
			goto WAIT
		default:
		}

		agentId := i + 2 // real agent uses ID=1; load test uses 2 ~ agentCount+1
		agent, err := NewSimAgent(agentId, *addr, metrics, *containers, *rateMs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "agent %d create failed: %v\n", agentId, err)
			continue
		}

		agentWg.Add(1)
		go func(a *SimAgent) {
			activeAgents.Add(1)
			defer activeAgents.Add(-1)
			a.Run(ctx, &agentWg)
		}(agent)

		started++

		if rampInterval > 0 {
			select {
			case <-ctx.Done():
				goto WAIT
			case <-time.After(rampInterval):
			}
		}
	}

	fmt.Printf("All %d agents started in %v\n\n", started, time.Since(startTime).Round(time.Millisecond))

	{
		reportTicker := time.NewTicker(10 * time.Second)
		defer reportTicker.Stop()

		var durationC <-chan time.Time
		if *duration > 0 {
			remaining := *duration - time.Since(startTime)
			if remaining <= 0 {
				cancel()
				goto WAIT
			}
			t := time.NewTimer(remaining)
			defer t.Stop()
			durationC = t.C
		}

		for {
			select {
			case <-ctx.Done():
				goto WAIT
			case <-durationC:
				fmt.Printf("Duration %v reached, shutting down...\n", *duration)
				cancel()
				goto WAIT
			case <-reportTicker.C:
				fmt.Print(metrics.Report(time.Since(startTime), int(activeAgents.Load())))
			}
		}
	}

WAIT:
	fmt.Println("Waiting for all agents to stop...")
	agentWg.Wait()
	fmt.Printf("\n%s", metrics.Summary(time.Since(startTime)))
}
