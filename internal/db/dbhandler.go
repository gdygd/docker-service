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
}
