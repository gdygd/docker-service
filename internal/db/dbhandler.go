package db

import (
	"context"
	"database/sql"
)

const (
	PENDING   = 1 // 대기
	CONFIRMED = 2 // 확정
	CANCELLED = 3 // 취소
)

type DbHandler interface {
	Init() error
	Close(*sql.DB)
	ReadSysdate(ctx context.Context) (string, error)
	ReadUser(ctx context.Context, username string) (User, error)
	ReadUserSession(ctx context.Context, id string) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateUserSession(ctx context.Context, arg CreateSessionParams) (Session, error)
}
