package service

import (
	"context"

	"docker_service/internal/db"
	"docker_service/internal/logger"
)

func (s *ApiService) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	user, err := s.dbHnd.CreateUser(ctx, arg)
	if err != nil {
		logger.Log.Error("[CreateUser] DB error: %v", err)
		return db.User{}, err
	}
	return user, nil
}