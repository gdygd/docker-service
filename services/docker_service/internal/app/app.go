package app

import (
	"docker_service/internal/container"
	"docker_service/internal/logger"
	"docker_service/internal/server/api"
	"docker_service/internal/server/pipe"
	"sync"
)

type Application struct {
	wg         *sync.WaitGroup
	ApiServer  *api.Server
	PipeServer *pipe.Server
}

func NewApplication(ct *container.Container, ch_terminate chan bool) *Application {
	var wg *sync.WaitGroup = &sync.WaitGroup{}

	// new httpserver
	apisvr, err := api.NewServer(wg, ct, ch_terminate)
	if err != nil {
		logger.Log.Error("Api server initialization fail.. %v", err)
		return nil
	}

	// Pipeline Server 초기화
	pipeCfg := pipe.Config{
		IntervalSec: 30, // 30초 주기
		BufferSize:  50,
	}
	pipesvr, err := pipe.NewServer(wg, ct.DockerMng, pipeCfg)
	if err != nil {
		logger.Log.Error("Pipe server initialization fail.. %v", err)
		return nil
	}

	return &Application{
		wg:         wg,
		ApiServer:  apisvr,
		PipeServer: pipesvr,
	}
}

func (app *Application) Start() {
	// API Server 시작
	app.wg.Add(1)
	logger.Log.Print(3, "Start API server..")
	go app.ApiServer.Start()

	// Pipeline Server 시작
	app.wg.Add(1)
	logger.Log.Print(3, "Start Pipe server..")
	go app.PipeServer.Start()
}

func (app *Application) Shutdown() {
	logger.Log.Print(3, "Shutdown Pipe server..")
	app.PipeServer.Shutdown()

	logger.Log.Print(3, "Shutdown API server..")
	app.ApiServer.Shutdown()

	logger.Log.Print(3, "Shutdown complete")
}

