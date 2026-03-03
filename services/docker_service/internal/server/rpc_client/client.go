package gapi

import (
	"context"
	"docker_service/internal/container"
	"docker_service/internal/logger"
	"docker_service/pb"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	wg     *sync.WaitGroup // app에서 관리하는 wg
	txrxwg sync.WaitGroup  // tx, rx routine group
	conn   *grpc.ClientConn
	stream pb.HelloService_ConnMessageClient
	ctx    context.Context // master context
	cancel context.CancelFunc
	mu     sync.RWMutex

	ct *container.Container
}

func NewClient(wg *sync.WaitGroup, ct *container.Container, ch_terminate chan bool) (*GrpcClient, error) {
	ctx, cancel := context.WithCancel(context.Background())
	gclient := &GrpcClient{
		wg:     wg,
		ctx:    ctx,
		cancel: cancel,
		ct:     ct,
	}

	err := gclient.Connect()
	if err != nil {
		return gclient, err
	}

	err = gclient.CreateStream()
	if err != nil {
		return gclient, err
	}

	return gclient, nil
}

func (c *GrpcClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	logger.Log.Print(2, "Connect...#1")

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient("localhost:9190", opts...)
	c.conn = conn
	if err != nil {
		logger.Log.Error("Failed to connect: %v", err)
		return err
	}

	logger.Log.Print(2, "Connect...#2")

	return nil
}

func (c *GrpcClient) CreateStream() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	logger.Log.Print(2, "CreateStream...#1")

	c.stream = nil

	if c.conn == nil {
		return fmt.Errorf("conn is nil..")
	}
	client := pb.NewHelloServiceClient(c.conn)
	stream, err := client.ConnMessage(context.Background())
	if err != nil {
		logger.Log.Error("Error creating stream.. %v", err)
		return err
	}

	c.stream = stream
	return nil
}

func (c *GrpcClient) Start() {
	defer c.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			log.Printf("THRclient panic: %v \n", r)
		}
	}()
	txrxCtx, txrxCancel := context.WithCancel(c.ctx) // master ctx로 하위 루틴 컨텍스트 생성

	var mu sync.Mutex
	var txRunning, rxRunning, connRunning bool // 실행상태 체크
	txDone := make(chan struct{}, 1)
	rxDone := make(chan struct{}, 1)
	connectDone := make(chan struct{}, 1)

	// 최초 시작
	c.startTxRoutine(txrxCtx, &mu, &txRunning, txDone)
	c.startRxRoutine(txrxCtx, &mu, &rxRunning, rxDone)
	c.startManageConnect(txrxCtx, &mu, &connRunning, connectDone)

	// check routine
	for {
		select {
		case <-c.ctx.Done():
			logger.Log.Print(2, "grpc client master context canceled..")
			txrxCancel()
			goto WAIT
		case <-txDone: // Manage txRoutine
			logger.Log.Warn("Tx routine ended..")
			if c.ctx.Err() == nil {
				c.startTxRoutine(txrxCtx, &mu, &txRunning, txDone)
			}
		case <-rxDone: // Manage rxRoutine
			logger.Log.Warn("Rx routine ended..")
			if c.ctx.Err() == nil {
				c.startRxRoutine(txrxCtx, &mu, &rxRunning, rxDone)
			}
		case <-connectDone:
			logger.Log.Warn("Connect routine ended..")
			if c.ctx.Err() == nil {
				c.startManageConnect(txrxCtx, &mu, &connRunning, connectDone)
			}
		default:
			time.Sleep(time.Second * 1)
		}
	}
WAIT:
	done := make(chan struct{})
	go func() {
		for {
			mu.Lock()
			if !txRunning && !rxRunning && !connRunning {
				mu.Unlock()
				break
			}
			mu.Unlock()
			time.Sleep(time.Second * 1)
		}
		close(connectDone)
	}()

	select {
	case <-done:
		log.Printf("Graceful shutdown !")
	case <-time.After(time.Second * 5):
		log.Printf("Graceful shutdown timeout.. client...")
	}

	logger.Log.Print(2, "gRPC Client quit..")
}

func (c *GrpcClient) Shutdown() {
	logger.Log.Print(2, "Shutdown grpc client..")

	// c.conn.Close()

	c.cancel() // 종료  시그널
	c.txrxwg.Wait()

	logger.Log.Print(2, "Tx/Rx routine finished.")
	// <-c.ctx.Done()
	logger.Log.Print(2, "gRPC client Shutdown finished.")
}
