package container

import (
	"api-gateway/internal/config"
	"api-gateway/internal/db"
	"api-gateway/internal/db/mdb"
	"api-gateway/internal/logger"
	"api-gateway/internal/memory"
	"fmt"
)

type Container struct {
	Config *config.Config
	DbHnd  db.DbHandler
	ObjDb  *memory.RedisDb
}

var (
	container *Container
	runMode   config.RunMode = config.ModeDev // 기본값: 개발 모드
)

// SetRunMode 실행 모드 설정 (container 생성 전에 호출)
func SetRunMode(mode int) {
	runMode = config.RunMode(mode)
}

func NewContainer() (*Container, error) {
	container = &Container{}
	// load config
	cfg, err := initConfig()
	if err != nil {
		return nil, fmt.Errorf("config loading error..%v \n", err)
	}
	container.Config = &cfg

	// init database
	dbhnd := initDatabase(cfg)
	container.DbHnd = dbhnd

	// init object db
	obj := memory.InitRedisDb(cfg.RedisAddr)
	container.ObjDb = obj

	return container, nil
}

func initConfig() (config.Config, error) {
	return config.LoadConfig(".", runMode)
}

func initDatabase(config config.Config) db.DbHandler {
	mdb := mdb.NewMdbHandler(config.DBUser, config.DBPasswd, config.DBSName, config.DBAddress, config.DBPort)
	err := mdb.Init()
	if err != nil {
		logger.Log.Error("Db Init err.. %v", err)
	}
	return mdb
}
