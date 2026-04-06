package container

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"docker_service/internal/config"
	"docker_service/internal/db"
	"docker_service/internal/db/mdb"
	"docker_service/internal/docker"
	"docker_service/internal/logger"
	"docker_service/internal/memory"

	"github.com/gdygd/goglib/databus"
)

type Container struct {
	Config    *config.Config
	DbHnd     db.DbHandler
	ObjDb     *memory.RedisDb
	Docker    *docker.Client
	DockerMng *docker.DockerClientManager

	Bus *databus.DataBus
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

	// init config from db
	strconfig := initHostInfo(dbhnd)
	container.Config.DOCKER_HOSTS = strconfig

	// init databus
	container.Bus = databus.NewDataBus()

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
			Mode: h.Mode,
		})
	}

	dockerMng, err := docker.NewDockerClientManager(dockerHosts)
	if err != nil {
		logger.Log.Error("init docker client error..(%v)", err)
	}
	container.DockerMng = dockerMng

	// // init databus
	// container.Bus = databus.NewDataBus()

	// strconfig := initHostInfo(dbhnd)
	// container.Config.DOCKER_HOSTS = strconfig

	return container, nil
}

func initConfig() (config.Config, error) {
	return config.LoadConfig(".")
}

func initDatabase(config config.Config) db.DbHandler {
	mdb := mdb.NewMdbHandler(config.DBUser, config.DBPasswd, config.DBSName, config.DBAddress, config.DBPort)

	for i := 0; i < 10; i++ {
		err := mdb.Init()
		if err != nil {
			logger.Log.Error("Db Init err.. (%d)%v", i, err)
		} else {
			logger.Log.Print(3, "Db init OK!")
			break
		}
		time.Sleep(2 * time.Second)
	}

	// err := mdb.Init()
	// if err != nil {
	// 	logger.Log.Error("Db Init err.. %v", err)
	// }
	return mdb
}

// DockerHostConfig는 Docker 호스트 설정
type DockerHostConfig struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Addr string `json:"addr"`
	Mode int    `json:"mode"` // 1:docker.sock, 2:tls
}

func initHostInfo(dbHnd db.DbHandler) string {
	hosts, _ := dbHnd.ReadHost(context.Background())
	hostinfo := []DockerHostConfig{}
	for _, host := range hosts {
		hostinfo = append(hostinfo, DockerHostConfig{
			Id:   host.HostId,
			Name: host.HostName,
			Addr: host.HostAddress,
			Mode: host.Mode,
		})
	}
	if len(hosts) == 0 {
		hostinfo = []DockerHostConfig{
			{
				Id:   1,
				Name: "Localhost",
				Addr: "unix:///var/run/docker.sock",
				Mode: 1,
			},
		}
	}
	data, _ := json.Marshal(hostinfo)

	fmt.Printf("##########%s \n", string(data))

	return string(data)
}
