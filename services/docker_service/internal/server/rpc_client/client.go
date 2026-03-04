package gapi

import (
	"context"
	"docker_service/internal/container"
	"docker_service/internal/logger"
	"docker_service/internal/pipeline"
	"docker_service/pb"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type GrpcClient struct {
	wg     *sync.WaitGroup
	conn   *grpc.ClientConn
	stream pb.ContainerService_DataStreamClient
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.RWMutex

	addr     string
	agentKey string
	ct       *container.Container
	pipeCh   <-chan pipeline.Message
}

func NewClient(
	wg *sync.WaitGroup,
	ct *container.Container,
	pipeCh <-chan pipeline.Message,
	addr string,
	agentKey string,
) (*GrpcClient, error) {
	ctx, cancel := context.WithCancel(context.Background())
	c := &GrpcClient{
		wg:       wg,
		ctx:      ctx,
		cancel:   cancel,
		ct:       ct,
		pipeCh:   pipeCh,
		addr:     addr,
		agentKey: agentKey,
	}

	if err := c.connect(); err != nil {
		cancel()
		return nil, err
	}
	if err := c.createStream(); err != nil {
		cancel()
		return nil, err
	}

	return c, nil
}

func (c *GrpcClient) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, err := grpc.NewClient(
		c.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(c.agentKeyUnaryInterceptor()),
		grpc.WithStreamInterceptor(c.agentKeyStreamInterceptor()),
	)
	if err != nil {
		return fmt.Errorf("grpc.NewClient(%s): %w", c.addr, err)
	}

	if c.conn != nil {
		_ = c.conn.Close()
	}
	c.conn = conn
	return nil
}

func (c *GrpcClient) agentKeyUnaryInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "agent-key", c.agentKey)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (c *GrpcClient) agentKeyStreamInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, "agent-key", c.agentKey)
		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (c *GrpcClient) createStream() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return fmt.Errorf("conn is nil")
	}

	stream, err := pb.NewContainerServiceClient(c.conn).DataStream(c.ctx)
	if err != nil {
		return fmt.Errorf("DataStream: %w", err)
	}

	c.stream = stream
	return nil
}

func (c *GrpcClient) resetStream() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.stream = nil
}

func (c *GrpcClient) getStream() pb.ContainerService_DataStreamClient {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.stream
}

func (c *GrpcClient) getConnState() connectivity.State {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.conn == nil {
		return connectivity.Shutdown
	}
	return c.conn.GetState()
}

func (c *GrpcClient) Start() {
	defer c.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("GrpcClient panic: %v", r)
		}
	}()

	txrxCtx, txrxCancel := context.WithCancel(c.ctx)
	defer txrxCancel()

	var stateMu sync.Mutex
	var txRunning, rxRunning, connRunning bool

	txDone := make(chan struct{}, 1)
	rxDone := make(chan struct{}, 1)
	connDone := make(chan struct{}, 1)

	launch := func(running *bool, done chan struct{}, name string, fn func(context.Context)) {
		stateMu.Lock()
		defer stateMu.Unlock()
		if *running {
			return
		}
		*running = true
		go func() {
			defer func() {
				stateMu.Lock()
				*running = false
				stateMu.Unlock()
				done <- struct{}{}
			}()
			logger.Log.Print(2, "[%s] started", name)
			fn(txrxCtx)
			logger.Log.Print(2, "[%s] exited", name)
		}()
	}

	launch(&connRunning, connDone, "manageConnect", c.manageConnect)
	launch(&txRunning, txDone, "txRoutine", c.txRoutine)
	launch(&rxRunning, rxDone, "rxRoutine", c.rxRoutine)

	for {
		select {
		case <-c.ctx.Done():
			logger.Log.Print(2, "gRPC client shutting down..")
			txrxCancel()
			goto WAIT

		case <-txDone:
			logger.Log.Warn("txRoutine exited")
			if c.ctx.Err() == nil {
				launch(&txRunning, txDone, "txRoutine", c.txRoutine)
			}

		case <-rxDone:
			logger.Log.Warn("rxRoutine exited")
			if c.ctx.Err() == nil {
				launch(&rxRunning, rxDone, "rxRoutine", c.rxRoutine)
			}

		case <-connDone:
			logger.Log.Warn("manageConnect exited")
			if c.ctx.Err() == nil {
				launch(&connRunning, connDone, "manageConnect", c.manageConnect)
			}
		}
	}

WAIT:
	stopped := make(chan struct{})
	go func() {
		for {
			stateMu.Lock()
			done := !txRunning && !rxRunning && !connRunning
			stateMu.Unlock()
			if done {
				close(stopped)
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	select {
	case <-stopped:
		logger.Log.Print(2, "Graceful shutdown complete")
	case <-time.After(5 * time.Second):
		logger.Log.Warn("Graceful shutdown timeout")
	}

	c.mu.Lock()
	if c.conn != nil {
		_ = c.conn.Close()
		c.conn = nil
	}
	c.mu.Unlock()

	logger.Log.Print(2, "gRPC Client quit")
}

func (c *GrpcClient) Shutdown() {
	logger.Log.Print(2, "Shutting down gRPC client..")
	c.cancel()
}
