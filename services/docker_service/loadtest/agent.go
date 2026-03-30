package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gdygd/goglib/databus"

	"docker_service/internal/container"
	"docker_service/internal/pipeline"
	gapi "docker_service/internal/server/rpc_client"
)

// SimAgent simulates a single agent: it generates synthetic pipeline.Message
// values at a fixed rate and forwards them to a GrpcClient for transmission.
type SimAgent struct {
	id       int
	pipeCh   chan pipeline.Message
	client   *gapi.GrpcClient
	gen      *Generator
	rateMs   time.Duration
	clientWg sync.WaitGroup
}

// NewSimAgent creates a SimAgent with agentId in range [2, 1001].
// It reuses the production GrpcClient unchanged; only the Container.Bus
// field is initialised since GrpcClient does not use any other field.
func NewSimAgent(id int, addr string, m *Metrics, containers, rateMs int) (*SimAgent, error) {
	pipeCh := make(chan pipeline.Message, 50)

	ct := &container.Container{
		Bus: databus.NewDataBus(),
	}

	a := &SimAgent{
		id:     id,
		pipeCh: pipeCh,
		gen:    NewGenerator(id, containers),
		rateMs: time.Duration(rateMs) * time.Millisecond,
	}

	a.clientWg.Add(1)
	client, err := gapi.NewClient(
		&a.clientWg,
		ct,
		pipeCh,
		addr,
		fmt.Sprintf("loadtest-%04d", id),
		gapi.WithUnaryInterceptor(m.Interceptor()),
	)
	if err != nil {
		return nil, fmt.Errorf("agent %d NewClient: %w", id, err)
	}
	a.client = client
	return a, nil
}

// Run starts the agent. It blocks until ctx is cancelled, then performs
// graceful shutdown of the GrpcClient.
func (a *SimAgent) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	// GrpcClient runs in its own goroutine and calls clientWg.Done on exit.
	go a.client.Start()

	ticker := time.NewTicker(a.rateMs)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			a.client.Shutdown()
			a.clientWg.Wait()
			return

		case <-ticker.C:
			msg := a.gen.Next()
			select {
			case a.pipeCh <- msg:
			default:
				// pipeCh is full; GrpcClient applies backpressure internally.
				// Drop this tick rather than blocking the generator.
			}
		}
	}
}
