package service

import (
	"auth-service/internal/db"
	"context"
)

type ServiceInterface interface {
	Test()

	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	LoginUser(ctx context.Context, username string) (db.User, error)
	CreateSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error)
	ReadSession(ctx context.Context, id string) (db.Session, error)
	DeleteSession(ctx context.Context, id string) error
}
