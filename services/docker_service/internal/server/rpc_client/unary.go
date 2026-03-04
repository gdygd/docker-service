package gapi

import (
	"context"
	"docker_service/pb"
	"fmt"
	"time"
)

// newServiceClient creates a stateless unary client per call.
func (c *GrpcClient) newServiceClient() (pb.ContainerServiceClient, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.conn == nil {
		return nil, fmt.Errorf("grpc connection is nil")
	}
	return pb.NewContainerServiceClient(c.conn), nil
}

func (c *GrpcClient) LoginUser(req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	client, err := c.newServiceClient()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()
	return client.LoginUser(ctx, req)
}

func (c *GrpcClient) ContainerState(req *pb.AgentMessage) (*pb.ServerMessage, error) {
	client, err := c.newServiceClient()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()
	return client.ContainerState(ctx, req)
}
