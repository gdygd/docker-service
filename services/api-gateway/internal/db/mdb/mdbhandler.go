package mdb

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewMdbHandler(user, pw, dbname, host string, port int) *MariaDbHandler {

	return &MariaDbHandler{user: user, pw: pw, dbNm: dbname, host: host, port: port}
}

type MariaDbHandler struct {
	db   *sql.DB
	user string
	pw   string
	dbNm string
	host string
	port int
	mu   sync.RWMutex
}

func (q *MariaDbHandler) Init() error {

	// dbSrc := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true&loc=Local", q.user, q.pw, q.host, q.dbNm)
	dbSrc := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", q.user, q.pw, q.host, q.port, q.dbNm)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	db, err := sql.Open("mysql", dbSrc)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(5)                  // 동시 최대 연결 수
	db.SetMaxIdleConns(3)                  // 유휴 상태로 유지할 연결 수
	db.SetConnMaxLifetime(1 * time.Minute) // 연결의 최대 수명

	// PingContext로 연결 확인
	if err := db.PingContext(ctx); err != nil {
		return err
	}

	fmt.Println("db open")
	q.mu.Lock()
	q.db = db
	q.mu.Unlock()

	return nil
}

func (q *MariaDbHandler) GetDB() *sql.DB {
	q.mu.RLock()
	defer q.mu.RUnlock()

	return q.db
}

func (q *MariaDbHandler) Close(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func (q *MariaDbHandler) execTx(ctx context.Context, fn func(*MariaDbHandler) error) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err : %v, rv err %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
