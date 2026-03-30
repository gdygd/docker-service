package app

import (
	"sync"

	"docker_service/internal/config"
	"docker_service/internal/container"
	evt "docker_service/internal/event2"
	"docker_service/internal/logger"
	"docker_service/internal/pipeline"
	"docker_service/internal/server/api"
	"docker_service/internal/server/event"
	"docker_service/internal/server/pipe"
	gapi "docker_service/internal/server/rpc_client"
)

type Application struct {
	wg         *sync.WaitGroup
	ApiServer  *api.Server
	PipeServer *pipe.Server
	Gclient    *gapi.GrpcClient
	config     *config.Config

	eventServer *event.Server
	pipeCh      chan pipeline.Message // pipe - rcp client간 데이터 전송 채널
}

func NewApplication(ct *container.Container, ch_terminate chan bool) *Application {
	var wg *sync.WaitGroup = &sync.WaitGroup{}

	pipeCh := make(chan pipeline.Message, 100)

	// event 수집 인스턴스 (evtMgr : container.Container 멤버로 관리 고려)
	evtMgr := evt.NewEventManager(ct.DockerMng)
	// event 수집 메니저 초기화
	evtsvr, err := event.NewServer(wg, ct, evtMgr) // evtMgr : watch host, and 이벤트 수집
	if err != nil {
		logger.Log.Error("Event server initialization fail.. %v", err)
		return nil
	}

	// new httpserver
	apisvr, err := api.NewServer(wg, ct, evtMgr) // evtMgr : 이벤트 구독
	if err != nil {
		logger.Log.Error("Api server initialization fail.. %v", err)
		return nil
	}

	if ct.Config.OprMode == "aws" {
		// Pipeline Server 초기화
		pipeCfg := pipe.Config{
			IntervalSec: 30, // 30초 주기
			BufferSize:  50,
		}
		pipesvr, err := pipe.NewServer(wg, ct.DockerMng, pipeCfg, ct.Config, pipeCh, evtMgr)
		if err != nil {
			logger.Log.Error("Pipe server initialization fail.. %v", err)
			return nil
		}

		// init grpc client
		// gclient, err := gapi.NewClient(wg, ct, pipeCh, "localhost:9190", "agentkey...")
		gclient, err := gapi.NewClient(wg, ct, pipeCh, ct.Config.AwsRpcServerAddress, "agentkey...")
		if err != nil {
			logger.Log.Error("gRPC client initialization fail.. %v", err)
			return nil
		}
		return &Application{
			wg:          wg,
			ApiServer:   apisvr,
			PipeServer:  pipesvr,
			Gclient:     gclient,
			pipeCh:      pipeCh,
			eventServer: evtsvr,
			config:      ct.Config,
		}
	}
	return &Application{
		wg:          wg,
		ApiServer:   apisvr,
		pipeCh:      pipeCh,
		eventServer: evtsvr,
		config:      ct.Config,
	}
}

func (app *Application) Start() {
	// API Server 시작
	app.wg.Add(1)
	logger.Log.Print(3, "Start API server..")
	go app.ApiServer.Start()

	if app.config.OprMode == "aws" {
		// Pipeline Server 시작
		app.wg.Add(1)
		logger.Log.Print(3, "Start Pipe server..")
		go app.PipeServer.Start()

		// gRPC client 시작
		app.wg.Add(1)
		logger.Log.Print(3, "Start gRPC client..")
		go app.Gclient.Start()

	}

	// event server 시작
	app.wg.Add(1)
	logger.Log.Print(3, "Start event server..")
	go app.eventServer.Start()
}

func (app *Application) Shutdown() {
	if app.config.OprMode == "aws" {
		logger.Log.Print(3, "Shutdown Pipe server..")
		app.PipeServer.Shutdown()
		close(app.pipeCh)

	}

	logger.Log.Print(3, "Shutdown API server..")
	app.ApiServer.Shutdown()

	if app.config.OprMode == "aws" {
		logger.Log.Print(3, "Shutdown grpc client..")
		go app.Gclient.Shutdown()
	}

	// event server 시작
	app.wg.Add(1)
	logger.Log.Print(3, "Shutdown event server..")
	go app.eventServer.Shutdown()

	logger.Log.Print(3, "Shutdown complete")
}
