package gapi

import (
	"context"

	"auth-service/internal/logger"

	"auth-service/pb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	logger.Log.Print(2, "name : %s, pw: %s", req.GetUsername(), req.GetPassword())
	return nil, nil
}
