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

func (s *ApiService) LoginUser(ctx context.Context, username string) (db.User, error) {
	user, err := s.dbHnd.ReadUser(ctx, username)
	if err != nil {
		logger.Log.Error("[LoginUser] DB error: %v", err)
		return db.User{}, err
	}
	return user, nil
}

func (s *ApiService) CreateSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error) {
	se, err := s.dbHnd.CreateUserSession(ctx, arg)
	if err != nil {
		logger.Log.Error("[CreateSession] DB error: %v", err)
		return db.Session{}, err
	}
	return se, nil
}

func (s *ApiService) ReadSession(ctx context.Context, id string) (db.Session, error) {
	user, err := s.dbHnd.ReadUserSession(ctx, id)
	if err != nil {
		logger.Log.Error("[ReadSession] DB error: %v", err)
		return db.Session{}, err
	}
	return user, nil
}

func (s *ApiService) DeleteSession(ctx context.Context, id string) error {
	err := s.dbHnd.DeleteUserSession(ctx, id)
	if err != nil {
		logger.Log.Error("[DeleteSession] DB error: %v", err)
		return err
	}
	return nil
}
