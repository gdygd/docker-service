package container

import (
	"docker_service/internal/config"
	"docker_service/internal/db"
	"docker_service/internal/db/mdb"
	"docker_service/internal/docker"
	"docker_service/internal/logger"
	"docker_service/internal/memory"
	"fmt"
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

	// init docker client (docker 소켓에 연결, local host)
	dcCli, err := docker.New()
	if err != nil {
		logger.Log.Error("init docker client error..(%v)", err)
	}
	container.Docker = dcCli

	// set docker daemon cert path
	docker.SetCertpaht(config.CERT_PATH)

	// init docker client Manager (설정 파일에서 호스트 목록 로드)
	hostConfigs, err := config.GetDockerHosts()
	if err != nil {
		logger.Log.Error("parse docker hosts config error..(%v)", err)
	}

	// config.DockerHostConfig -> docker.HostConfig 변환
	dockerHosts := make([]docker.HostConfig, 0, len(hostConfigs))
	for _, h := range hostConfigs {
		dockerHosts = append(dockerHosts, docker.HostConfig{
			Name: h.Name,
			Addr: h.Addr,
		})
	}

	dockerMng, err := docker.NewDockerClientManager(dockerHosts)
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
	err := mdb.Init()
	if err != nil {
		logger.Log.Error("Db Init err.. %v", err)
	}
	return mdb
}
