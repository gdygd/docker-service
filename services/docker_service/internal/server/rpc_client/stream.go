package gapi

import (
	"context"
	"docker_service/internal/logger"
	"docker_service/pb"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/gdygd/goglib/databus"
	"google.golang.org/grpc/connectivity"
)

func (c *GrpcClient) startManageConnect(ctx context.Context, mu *sync.Mutex, running *bool, done chan struct{}) {
	mu.Lock()
	if *running {
		mu.Unlock()
		return
	}
	*running = true
	mu.Unlock()

	c.txrxwg.Add(1)
	go func() {
		defer c.txrxwg.Done()
		defer func() {
			mu.Lock()
			*running = false
			mu.Unlock()
			done <- struct{}{}
		}()
		log.Printf("Starting Manage Connect routine for client")
		c.manageConnect(ctx)
	}()
}

func (c *GrpcClient) startTxRoutine(ctx context.Context, mu *sync.Mutex, running *bool, done chan struct{}) {
	mu.Lock()
	if *running {
		mu.Unlock()
		return
	}
	*running = true
	mu.Unlock()

	c.txrxwg.Add(1)
	go func() {
		defer c.txrxwg.Done()
		defer func() {
			mu.Lock()
			*running = false
			mu.Unlock()
			done <- struct{}{}
		}()
		log.Printf("Starting Tx routine for client")
		c.txRoutine(ctx)
	}()
}

func (c *GrpcClient) startRxRoutine(ctx context.Context, mu *sync.Mutex, running *bool, done chan struct{}) {
	mu.Lock()
	if *running {
		mu.Unlock()
		return
	}
	*running = true
	mu.Unlock()

	c.txrxwg.Add(1)
	go func() {
		defer c.txrxwg.Done()
		defer func() {
			mu.Lock()
			*running = false
			mu.Unlock()
			done <- struct{}{}
		}()
		log.Printf("Starting Rx routine for client %d")
		c.rxRoutine(ctx)
	}()
}

func (c *GrpcClient) manageConnect(ctx context.Context) {
	for {
		logger.Log.Print(2, "manageconnect...#0")
		select {
		case <-ctx.Done():
			logger.Log.Print(2, "[ManageConn] routine exiting..")
			return
		default:
			logger.Log.Print(2, "manageconnect...#1 %v", c.getState())

			if c.getState() == connectivity.Shutdown || c.getState() == connectivity.TransientFailure || c.getState() == connectivity.Idle {
				logger.Log.Print(2, "manageconnect...#2")

				err1 := c.Connect()
				err2 := c.CreateStream()

				if err1 == nil && err2 == nil {
					logger.Log.Print(2, "Reconnect Successfully")
				} else {
					logger.Log.Error("gRPC Reconnect failed..[%v][%v]", err1, err2)
				}
			}

			time.Sleep(time.Second * 1)
		}
	}
}

func (c *GrpcClient) txRoutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Print(2, "quit tx routine...")
			if c.stream == nil {
				logger.Log.Print(2, "#1 tx, stream is nil..")
				return
			}

			// 남은 메세지 처리

			// 종료
			c.closeSend()
			return

		default:
			logger.Log.Print(2, "txRoutine... %v", c.getState())
			if c.getState() != connectivity.Ready {
				time.Sleep(time.Second * 1)
				continue
			}

			if c.stream == nil {
				logger.Log.Print(2, "#2 tx, stream is nil..")
				continue
			}

			// send.
			// c.mu.Lock()
			// err := c.stream.Send(&pb.Hello{Msg: "to Client..."})
			// c.mu.Unlock()
			err := c.send(&pb.Hello{Msg: "to Client..."})
			if err != nil {
				log.Printf("Send error: %v", err)
				return
			}

			time.Sleep(time.Second * 1)
		}
	}
}

func (c *GrpcClient) rxRoutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Print(2, "quit rx routine...")
			return
		default:
			if c.getState() != connectivity.Ready {
				continue
			}

			// c.mu.Lock()
			logger.Log.Print(2, "rxRoutine... %v", c.getState())
			// c.mu.Unlock()

			if c.getState() != connectivity.Ready {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			if c.stream == nil {
				logger.Log.Print(1, "rx, stream is nil..")
				continue
			}

			// recv
			resp, err := c.recv()

			// c.mu.Lock()
			// resp, err := c.stream.Recv()
			// c.mu.Unlock()

			// if err == io.EOF {
			// 	logger.Log.Error("Server closed stream")
			// 	return
			// }
			// if err != nil {
			// 	logger.Log.Error("Recv error: %v", err)
			// 	return
			// }
			if err == nil {
				log.Printf("From server: %s", resp.Msg)
				msg := databus.Message{
					Topic: "test_msg",
					Data:  resp.Msg,
				}
				c.ct.Bus.Publish(msg)

			}
		}
	}
}

func (c *GrpcClient) send(data interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	err := c.stream.Send(&pb.Hello{Msg: "to Client..."})
	if err != nil {
		logger.Log.Error("grpc Send err.. %v", err)
	}

	return err
}

func (c *GrpcClient) recv() (*pb.Hello, error) {
	c.mu.RLock()
	logger.Log.Print(2, "recv #1")
	defer func() {
		logger.Log.Print(2, "recv #2")
		c.mu.RUnlock()
	}()
	// defer c.mu.RUnlock()
	resp, err := c.stream.Recv()

	if resp == nil {
		return nil, fmt.Errorf("resp is nil...")
	}

	if err == io.EOF {
		logger.Log.Error("Server closed stream")
		return resp, nil
	}
	if err != nil {
		logger.Log.Error("Recv error: %v", err)
		return resp, nil
	}

	return resp, nil
}

func (c *GrpcClient) closeSend() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.stream.CloseSend()
	if err != nil {
		logger.Log.Error("closeSend err.. %v", err)
	}
	return err
}

func (c *GrpcClient) getState() connectivity.State {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn.GetState()
}
