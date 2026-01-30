package db

import (
	"context"
	"database/sql"
)

type DbHandler interface {
	Init() error
	Close(*sql.DB)
	ReadSysdate(ctx context.Context) (string, error)
	ReadUser(ctx context.Context, username string) (User, error)
	ReadUserSession(ctx context.Context, id string) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateUserSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	DeleteUserSession(ctx context.Context, id string) error
}
