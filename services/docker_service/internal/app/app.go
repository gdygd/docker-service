package app

import (
	"docker_service/internal/container"
	"docker_service/internal/logger"
	"docker_service/internal/pipeline"
	"docker_service/internal/server/api"
	"docker_service/internal/server/pipe"
	gapi "docker_service/internal/server/rpc_client"
	"sync"
)

type Application struct {
	wg         *sync.WaitGroup
	ApiServer  *api.Server
	PipeServer *pipe.Server
	Gclient    *gapi.GrpcClient

	pipeCh chan pipeline.Message // pipe - rcp client간 데이터 전송 채널
}

func NewApplication(ct *container.Container, ch_terminate chan bool) *Application {
	var wg *sync.WaitGroup = &sync.WaitGroup{}

	pipeCh := make(chan pipeline.Message, 100)

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
	pipesvr, err := pipe.NewServer(wg, ct.DockerMng, pipeCfg, pipeCh)
	if err != nil {
		logger.Log.Error("Pipe server initialization fail.. %v", err)
		return nil
	}

	// init grpc client
	// gclient, _ := gapi.NewClient(wg, ct, ch_terminate, pipeCh)
	gclient, _ := gapi.NewClient(wg, ct, pipeCh, "localhost:9190", "agentkey...")

	return &Application{
		wg:         wg,
		ApiServer:  apisvr,
		PipeServer: pipesvr,
		Gclient:    gclient,
		pipeCh:     pipeCh,
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

	// gRPC client 시작
	app.wg.Add(1)
	logger.Log.Print(3, "Start gRPC client..")
	go app.Gclient.Start()
}

func (app *Application) Shutdown() {
	logger.Log.Print(3, "Shutdown Pipe server..")
	app.PipeServer.Shutdown()
	close(app.pipeCh)

	logger.Log.Print(3, "Shutdown API server..")
	app.ApiServer.Shutdown()

	logger.Log.Print(3, "Shutdown grpc client..")
	go app.Gclient.Shutdown()

	logger.Log.Print(3, "Shutdown complete")
}
