package memory

import (
	"context"
	"sync"

	"docker_service/internal/cache"
)

type RedisDb struct {
	ctx    context.Context
	Rdb    *cache.RedisClient
	User   map[string]User
	SysStt SystemInfo
	usermu sync.Mutex
	sysmu  sync.Mutex
}

func InitRedisDb(addr string) *RedisDb {
	// init database

	RedisDb := &RedisDb{
		ctx:    context.Background(),
		Rdb:    cache.NewRedisClient(addr),
		User:   make(map[string]User),
		SysStt: SystemInfo{SvrUtc: 0, DbSvrComm: STATE_OFF, RdSvrComm: STATE_OFF},
	}

	return RedisDb
}
