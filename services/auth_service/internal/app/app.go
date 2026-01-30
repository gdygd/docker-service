package app

import (
	"sync"

	"auth-service/internal/container"
	"auth-service/internal/logger"
	"auth-service/internal/server/api"
	"auth-service/internal/server/gapi"
)

type Application struct {
	wg           *sync.WaitGroup
	ApiServer    *api.Server
	GApiServer   *gapi.Server
	GrpcGwServer *gapi.Server
}

func NewApplication(ct *container.Container, ch_terminate chan bool) *Application {
	var wg *sync.WaitGroup = &sync.WaitGroup{}

	// new httpserver
	// newserver(container)
	apisvr, err := api.NewServer(wg, ct, ch_terminate)
	if err != nil {
		logger.Log.Error("Api server initialization fail.. %v", err)
		return nil
	}

	gapisvr, err := gapi.NewServer(wg, ct, ch_terminate)
	if err != nil {
		logger.Log.Error("gRPC server initialization fail.. %v", err)
		return nil
	}

	grpcGwsvr, err := gapi.NewGatewayServer(wg, ct, ch_terminate)
	if err != nil {
		logger.Log.Error("gRPC GW server initialization fail.. %v", err)
		return nil
	}

	return &Application{
		wg:           wg,
		ApiServer:    apisvr,
		GApiServer:   gapisvr,
		GrpcGwServer: grpcGwsvr,
	}
}

func (app Application) Start() {
	app.wg.Add(1)
	logger.Log.Print(3, "Start API server.. #1")
	go app.ApiServer.Start()

	app.wg.Add(1)
	logger.Log.Print(3, "Start gRPC server.. #1")
	go app.GApiServer.StartgPRC()

	app.wg.Add(1)
	logger.Log.Print(3, "Start gRPC-GW server.. #1")
	go app.GrpcGwServer.StartgRPCGateway()
}

func (app Application) Shutdown() {
	logger.Log.Print(3, "Shutdown Rest server#1")
	app.ApiServer.Shutdown()
	logger.Log.Print(3, "Shutdown Rest server#2")

	logger.Log.Print(3, "Shutdown gRPC server#1")
	app.GApiServer.ShutdowngRPC()
	logger.Log.Print(3, "Shutdown gRPC server#2")

	logger.Log.Print(3, "Shutdown gRPC-GW server#1")
	app.GrpcGwServer.ShutdowngRPCGw()
	logger.Log.Print(3, "Shutdown gRPC-GW server#2")
}
