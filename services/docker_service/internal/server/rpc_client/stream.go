package gapi

import (
	"context"
	"docker_service/internal/logger"
	"docker_service/internal/pipeline"
	"docker_service/pb"
	"io"
	"time"

	"github.com/gdygd/goglib/databus"
	"google.golang.org/grpc/connectivity"
)

// manageConnect: 커넥션/스트림 상태 감시 및 복구
// 서버가 내려가 있어도 주기적으로 재시도하며, 올라오면 자동으로 스트림을 생성한다.
func (c *GrpcClient) manageConnect(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Print(2, "[manageConnect] exiting")
			return
		default:
		}

		// ClientConn이 명시적으로 Shutdown된 경우에만 재생성
		// TransientFailure는 gRPC 내부 backoff가 처리하므로 개입하지 않는다.
		if c.getConnState() == connectivity.Shutdown {
			logger.Log.Warn("[manageConnect] connection shutdown, recreating ClientConn..")
			if err := c.connect(); err != nil {
				logger.Log.Error("[manageConnect] connect failed: %v", err)
				time.Sleep(3 * time.Second)
				continue
			}
		}

		// stream이 없으면 서버 상태와 무관하게 생성 시도
		// 서버가 내려가 있으면 실패하고 재시도, 올라오면 성공한다.
		if c.getStream() == nil {
			logger.Log.Print(2, "[manageConnect] no stream, attempting createStream..")
			if err := c.createStream(); err != nil {
				logger.Log.Warn("[manageConnect] createStream failed (retry in 3s): %v", err)
				time.Sleep(3 * time.Second)
				continue
			}
			logger.Log.Print(2, "[manageConnect] stream created")
		}

		time.Sleep(time.Second)
	}
}

// txRoutine: pipeCh에서 pipeline.Message를 읽어 타입에 따라 스트리밍/단항 RPC로 라우팅
func (c *GrpcClient) txRoutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Print(2, "[txRoutine] exiting")
			c.closeSend()
			return

		case msg, ok := <-c.pipeCh:
			if !ok {
				logger.Log.Warn("[txRoutine] pipeCh closed, initiating shutdown")
				c.cancel()
				return
			}

			switch msg.Type {
			// case pipeline.DataTypeList, pipeline.DataTypeStats, pipeline.DataTypeEvent:
			case pipeline.DataTypeList:
				c.sendStream(msg) // 실시간 스트리밍
			// case pipeline.DataTypeInspect:
			// 	c.sendUnary(msg) // 단발성 스냅샷
			default:
				logger.Log.Warn("[txRoutine] unknown message type: %s", msg.Type)
			}
		}
	}
}

// sendStream: DataStream 양방향 스트림으로 전송 (List, Stats, Event)
func (c *GrpcClient) sendStream(msg pipeline.Message) {
	logger.Log.Print(2, "sendStream...")
	if c.getConnState() != connectivity.Ready {
		logger.Log.Warn("[sendStream] connection not ready, dropping %s", msg.Type)
		return
	}

	stream := c.getStream()
	if stream == nil {
		logger.Log.Warn("[sendStream] stream is nil, dropping %s", msg.Type)
		return
	}

	pbMsg := &pb.AgentMessage{
		AgentKey: "abcd",
		Host:     "yun119",
	}

	// pbMsg, err := ConvertToAgentMessage(msg, c.agentKey)
	// if err != nil {
	// 	logger.Log.Error("[sendStream] convert failed: %v", err)
	// 	return
	// }

	if err := stream.Send(pbMsg); err != nil {
		logger.Log.Error("[sendStream] Send error: %v", err)
		c.resetStream()
	}
}

// sendUnary: ContainerState 단항 RPC 호출 (Inspect 스냅샷)
func (c *GrpcClient) sendUnary(msg pipeline.Message) {
	pbMsg, err := ConvertToAgentMessage(msg, c.agentKey)
	if err != nil {
		logger.Log.Error("[sendUnary] convert failed: %v", err)
		return
	}

	resp, err := c.ContainerState(pbMsg)
	if err != nil {
		logger.Log.Error("[sendUnary] ContainerState error: %v", err)
		return
	}

	c.handleServerMessage(resp)
}

// rxRoutine: DataStream에서 ServerMessage를 수신하여 databus에 발행
func (c *GrpcClient) rxRoutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Print(2, "[rxRoutine] exiting")
			return
		default:
		}

		if c.getConnState() != connectivity.Ready {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		stream := c.getStream()
		if stream == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		resp, err := stream.Recv()
		if err == io.EOF {
			logger.Log.Warn("[rxRoutine] server closed stream")
			c.resetStream()
			return
		}
		if err != nil {
			logger.Log.Error("[rxRoutine] Recv error: %v", err)
			c.resetStream()
			return
		}

		logger.Log.Print(2, "recv : %v", resp)
		// c.handleServerMessage(resp)
	}
}

func (c *GrpcClient) handleServerMessage(msg *pb.ServerMessage) {
	c.ct.Bus.Publish(databus.Message{
		Topic: "server_command",
		Data:  msg,
	})
}

func (c *GrpcClient) closeSend() {
	stream := c.getStream()
	if stream == nil {
		return
	}
	if err := stream.CloseSend(); err != nil {
		logger.Log.Error("[closeSend] error: %v", err)
	}
}
