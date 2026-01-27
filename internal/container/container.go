package container

import (
	"fmt"

	"docker_service/internal/config"
	"docker_service/internal/db"
	"docker_service/internal/db/mdb"
	"docker_service/internal/docker"
	"docker_service/internal/logger"
	"docker_service/internal/memory"
)

type Container struct {
	Config    *config.Config
	DbHnd     db.DbHandler
	ObjDb     *memory.RedisDb
	Docker    *docker.Client
	DockerMng *docker.DockerClientManager
}

var container *Container

func NewContainer() (*Container, error) {
	container = &Container{}
	// load config
	config, err := initConfig()
	if err != nil {
		return nil, fmt.Errorf("config loading error..%v \n", err)
	}
	container.Config = &config

	// init database
	dbhnd := initDatabase(config)
	container.DbHnd = dbhnd

	// init object db
	obj := memory.InitRedisDb(config.RedisAddr)
	container.ObjDb = obj

	// init docker client
	dcCli, err := docker.New()
	if err != nil {
		logger.Log.Error("init docker client error..(%v)", err)
	}
	container.Docker = dcCli

	// set docker daemon cert path
	docker.SetCertpaht(config.CERT_PATH)

	// init docker client Manager
	dockerMng, err := docker.NewDockerClientManager([]docker.HostConfig{
		// {Name: "local", Addr: "unix"},
		{Name: "119server", Addr: "tcp://10.1.0.119:2376"},
		// {Name: "localhost", Addr: "tcp://10.1.0.119:2375"}, // ssh tunnel 로연결
	})
	if err != nil {
		logger.Log.Error("init docker client error..(%v)", err)
	}
	container.DockerMng = dockerMng

	return container, nil
}

func initConfig() (config.Config, error) {
	return config.LoadConfig(".")
}

func initDatabase(config config.Config) db.DbHandler {
	mdb := mdb.NewMdbHandler(config.DBUser, config.DBPasswd, config.DBSName, config.DBAddress, config.DBPort)
	mdb.Init()
	return mdb
}
