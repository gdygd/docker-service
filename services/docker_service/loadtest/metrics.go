package main

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
)

const (
	idxList    = 0
	idxInspect = 1
	idxStats   = 2
	idxEvent   = 3
	idxCount   = 4
)

var typeNames = [idxCount]string{"list", "inspect", "stats", "event"}

// Metrics collects per-RPC-type success/error counts and cumulative latency.
// All fields are updated atomically from concurrent goroutines.
type Metrics struct {
	success    [idxCount]atomic.Int64
	errors     [idxCount]atomic.Int64
	latencyNs  [idxCount]atomic.Int64 // cumulative nanoseconds
	latencyCnt [idxCount]atomic.Int64 // sample count
}

func (m *Metrics) record(idx int, lat time.Duration, err error) {
	if err != nil {
		m.errors[idx].Add(1)
		return
	}
	m.success[idx].Add(1)
	m.latencyNs[idx].Add(lat.Nanoseconds())
	m.latencyCnt[idx].Add(1)
}

// Interceptor returns a gRPC unary client interceptor that records metrics.
func (m *Metrics) Interceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		idx := methodToIdx(method)
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		m.record(idx, time.Since(start), err)
		return err
	}
}

// methodToIdx maps gRPC full method name to a DataType index.
func methodToIdx(method string) int {
	switch {
	case strings.HasSuffix(method, "ContainerInfo"):
		return idxList
	case strings.HasSuffix(method, "ContainerInspect"):
		return idxInspect
	case strings.HasSuffix(method, "ContainerStats"):
		return idxStats
	case strings.HasSuffix(method, "ContainerEvent"):
		return idxEvent
	default:
		return idxList
	}
}

// Report returns a formatted metrics snapshot.
func (m *Metrics) Report(elapsed time.Duration, activeAgents int) string {
	var totalOk, totalErr int64
	for i := 0; i < idxCount; i++ {
		totalOk += m.success[i].Load()
		totalErr += m.errors[i].Load()
	}

	tps := 0.0
	if s := elapsed.Seconds(); s > 0 {
		tps = float64(totalOk) / s
	}
	errRate := 0.0
	if total := totalOk + totalErr; total > 0 {
		errRate = float64(totalErr) / float64(total) * 100
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "[T+%ds] agents=%d tps=%.1f err_rate=%.2f%%\n",
		int(elapsed.Seconds()), activeAgents, tps, errRate)

	for i := 0; i < idxCount; i++ {
		ok := m.success[i].Load()
		errCnt := m.errors[i].Load()
		avgMs := 0.0
		if cnt := m.latencyCnt[i].Load(); cnt > 0 {
			avgMs = float64(m.latencyNs[i].Load()) / float64(cnt) / 1e6
		}
		fmt.Fprintf(&sb, "  %-8s ok=%-8d err=%-5d avg=%.2fms\n",
			typeNames[i], ok, errCnt, avgMs)
	}
	return sb.String()
}

// Summary returns final aggregated report.
func (m *Metrics) Summary(elapsed time.Duration) string {
	return "[FINAL SUMMARY]\n" + m.Report(elapsed, 0)
}
