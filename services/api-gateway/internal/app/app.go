package app

import (
	"sync"

	"api-gateway/internal/container"
	"api-gateway/internal/logger"
	"api-gateway/internal/server/api"
)

type Application struct {
	wg        *sync.WaitGroup
	ApiServer *api.Server
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

	return &Application{
		wg:        wg,
		ApiServer: apisvr,
	}
}

func (app Application) Start() {
	app.wg.Add(1)
	logger.Log.Print(3, "Start API-GW server.. #1")
	go app.ApiServer.Start()
}

func (app Application) Shutdown() {
	logger.Log.Print(3, "Shutdown API-GW server#1")
	app.ApiServer.Shutdown()
	logger.Log.Print(3, "Shutdown API-GW server#2")
	// app.Terminate <- true
}
