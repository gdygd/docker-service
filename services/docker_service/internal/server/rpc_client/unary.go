package gapi

import (
	"context"
	"docker_service/pb"
	"time"
)

// unaryClient 매번 새로 생성 (stateless)
func (c *GrpcClient) newServiceClient() pb.ContainerServiceClient {
	return pb.NewContainerServiceClient(c.conn)
}

// LoginUser
func (c *GrpcClient) LoginUser(req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()

	return c.newServiceClient().LoginUser(ctx, req)
}

func (c *GrpcClient) ContainerState(req *pb.AgentMessage) (*pb.ServerMessage, error) {
	ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
	defer cancel()

	return c.newServiceClient().ContainerState(ctx, req)
}
