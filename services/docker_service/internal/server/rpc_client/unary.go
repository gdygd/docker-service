package gapi

import (
	"context"
	"fmt"
	"time"

	"docker_service/internal/logger"
	"docker_service/pb"
)

// newServiceClient creates a stateless unary client per call.
func (c *GrpcClient) newServiceClient() (pb.ContainerServiceClient, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.conn == nil {
		logger.Log.Error("grpc connection is nil")
		return nil, fmt.Errorf("grpc connection is nil")
	}
	return pb.NewContainerServiceClient(c.conn), nil
}

func (c *GrpcClient) LoginUser(req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	client, err := c.newServiceClient()
	if err != nil {
		logger.Log.Error("LoginUser error : %v", err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()
	return client.LoginUser(ctx, req)
}

func (c *GrpcClient) ContainerState(req *pb.AgentMessage) (*pb.ServerMessage, error) {
	client, err := c.newServiceClient()
	if err != nil {
		logger.Log.Error("ContainerState error : %v", err)
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()
	return client.ContainerState(ctx, req)
}

func (c *GrpcClient) ContainerInfo(req *pb.AgentMessage) (*pb.ServerMessage, error) {
	logger.Log.Print(1, "ContainerInfo..")
	client, err := c.newServiceClient()
	if err != nil {
		logger.Log.Error("ContainerState error : %v", err)
		return nil, fmt.Errorf("[ContainerInfo] error : %w", err)
	}
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()
	return client.ContainerInfo(ctx, req)
}

func (c *GrpcClient) ContainerInspect(req *pb.AgentMessage) (*pb.ServerMessage, error) {
	logger.Log.Print(1, "ContainerInspect..")
	client, err := c.newServiceClient()
	if err != nil {
		logger.Log.Error("ContainerInspect error : %v", err)
		return nil, fmt.Errorf("[ContainerInspect] error : %w", err)
	}
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()
	return client.ContainerInspect(ctx, req)
}

func (c *GrpcClient) ContainerStats(req *pb.AgentMessage) (*pb.ServerMessage, error) {
	logger.Log.Print(1, "ContainerStats..")
	client, err := c.newServiceClient()
	if err != nil {
		logger.Log.Error("ContainerStats error : %v", err)
		return nil, fmt.Errorf("[ContainerStats] error : %w", err)
	}
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()
	return client.ContainerStats(ctx, req)
}

func (c *GrpcClient) ContainerEvent(req *pb.AgentMessage) (*pb.ServerMessage, error) {
	logger.Log.Print(1, "ContainerEvent..")
	client, err := c.newServiceClient()
	if err != nil {
		logger.Log.Error("ContainerEvent error : %v", err)
		return nil, fmt.Errorf("[ContainerEvent] error : %w", err)
	}
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()
	return client.ContainerEvent(ctx, req)
}
