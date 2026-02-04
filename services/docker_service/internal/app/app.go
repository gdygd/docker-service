package app

import (
	"context"
	"docker_service/internal/container"
	"docker_service/internal/logger"
	"docker_service/internal/pipeline"
	"docker_service/internal/pipeline/collector"
	"docker_service/internal/server/api"
	"sync"
)

type Application struct {
	wg             *sync.WaitGroup
	ApiServer      *api.Server
	PipelineMgr    *collector.Manager
	pipelineCtx    context.Context
	pipelineCancel context.CancelFunc
}

func NewApplication(ct *container.Container, ch_terminate chan bool) *Application {
	var wg *sync.WaitGroup = &sync.WaitGroup{}

	// new httpserver
	apisvr, err := api.NewServer(wg, ct, ch_terminate)
	if err != nil {
		logger.Log.Error("Api server initialization fail.. %v", err)
		return nil
	}

	// Pipeline Manager 초기화
	pipelineMgr := collector.NewManager(ct.DockerMng, 100)

	// 모든 호스트에 List Collector 등록
	cfg := collector.Config{
		IntervalSec: 30, // 30초 주기
		BufferSize:  50,
	}
	if err := pipelineMgr.RegisterAllHosts([]collector.CollectorType{collector.TypeList}, cfg); err != nil {
		logger.Log.Error("Pipeline collector registration fail.. %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Application{
		wg:             wg,
		ApiServer:      apisvr,
		PipelineMgr:    pipelineMgr,
		pipelineCtx:    ctx,
		pipelineCancel: cancel,
	}
}

func (app *Application) Start() {
	// API Server 시작
	app.wg.Add(1)
	logger.Log.Print(3, "Start API server..")
	go app.ApiServer.Start()

	// Pipeline 시작
	app.wg.Add(1)
	go app.startPipeline()
}

func (app *Application) startPipeline() {
	defer app.wg.Done()

	outCh, err := app.PipelineMgr.Start(app.pipelineCtx)
	if err != nil {
		logger.Log.Error("Pipeline start fail.. %v", err)
		return
	}

	logger.Log.Print(3, "Pipeline started, collectors: %d", app.PipelineMgr.GetCollectorCount())

	// 수집된 메시지 처리
	for msg := range outCh {
		// TODO: gRPC Sender로 전송
		logger.Log.Print(2, "[Pipeline] type=%s host=%s timestamp=%v",
			msg.Type, msg.Host, msg.Timestamp)

		if msg.Type == "list" {
			containers := msg.Data.(pipeline.ContainerListData)
			logger.Log.Print(2, "Container List >> ")
			for _, c := range containers.Containers {
				logger.Log.Print(2, "\t ID:%s, Name:%s, Image:%s, State:%s, Status:%s ",
					c.ID, c.Name, c.Image, c.State, c.Status)
			}

		}

	}
}

func (app *Application) Shutdown() {
	logger.Log.Print(3, "Shutdown Pipeline..")
	app.pipelineCancel()
	app.PipelineMgr.Stop()

	logger.Log.Print(3, "Shutdown API server..")
	app.ApiServer.Shutdown()

	logger.Log.Print(3, "Shutdown complete")
}

/*
pipeline shutdown
┌─────────────────────────────────────────────────────────────────┐
│  1. Shutdown() 호출                                              │
│     app.pipelineCancel()  ← context 취소                         │
│     app.PipelineMgr.Stop()                                       │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  2. Manager.Stop()                                               │
│     - 각 Collector의 Stop() 호출                                 │
│     - close(m.outCh)  ← 출력 채널 닫음                           │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│  3. startPipeline()                                              │
│     for msg := range outCh  ← 채널이 닫히면 루프 종료            │
│     defer app.wg.Done()     ← WaitGroup 완료                     │
└─────────────────────────────────────────────────────────────────┘

*/
