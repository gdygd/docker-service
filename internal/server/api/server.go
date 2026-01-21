package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"docker_service/internal/config"
	"docker_service/internal/container"
	"docker_service/internal/db"
	"docker_service/internal/docker"
	"docker_service/internal/logger"
	"docker_service/internal/server/ws"
	"docker_service/internal/service"
	apiserv "docker_service/internal/service/api"

	"github.com/gdygd/goglib/token"

	"github.com/gin-gonic/gin"
)

const (
	R_TIME_OUT = 5 * time.Second
	W_TIME_OUT = 5 * time.Second
)

// Server serves HTTP requests for our banking service.
type Server struct {
	ctx context.Context

	wg           *sync.WaitGroup
	srv          *http.Server
	config       *config.Config
	tokenMaker   token.Maker
	router       *gin.Engine
	hub          *ws.Hub
	svr_cancel   context.CancelFunc
	service      service.ServiceInterface
	dbHnd        db.DbHandler
	ch_terminate chan bool
}

func NewServer(wg *sync.WaitGroup, ct *container.Container, ch_terminate chan bool) (*Server, error) {
	// init service
	apiservice := apiserv.NewApiService(ct.DbHnd, ct.Docker)
	tokenMaker, err := token.NewJWTMaker(ct.Config.TokenSecretKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker:%w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	server := &Server{
		ctx:          ctx,
		wg:           wg,
		config:       ct.Config,
		tokenMaker:   tokenMaker,
		service:      apiservice,
		dbHnd:        ct.DbHnd,
		ch_terminate: ch_terminate,
		hub:          ws.NewHub(ctx),
		svr_cancel:   cancel,
	}

	server.setupRouter()

	server.srv = &http.Server{}
	server.srv.Addr = ct.Config.HTTPServerAddress
	server.srv.Handler = server.router.Handler()
	server.srv.ReadTimeout = R_TIME_OUT
	server.srv.WriteTimeout = W_TIME_OUT

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	// gin.SetMode(gin.DebugMode)
	fmt.Printf("%v, \n", server.config.AllowOrigins)

	addresses := strings.Split(server.config.AllowOrigins, ",")

	router.GET("/heartbeat", server.heartbeat)
	router.GET("/terminate", server.terminate)

	router.GET("/ps", server.dockerPs)
	router.GET("/inspect/:id", server.containerInspect)
	router.GET("/start/:id", server.startContainer)
	router.GET("/stop/:id", server.stopContainer)
	router.GET("/stat/:id", server.statContainer)

	router.GET("/ws", server.wsHandler)

	// build
	// push
	// run
	// stop
	// start
	// restart
	// rm

	router.Use(corsMiddleware(addresses))
	router.Use(authMiddleware(server.tokenMaker))

	router.GET("/test", server.testapi)

	server.router = router
}

func (server *Server) updateContainerStats() {
	ticker := time.NewTicker(1 * time.Millisecond * 100)
	defer ticker.Stop()

	for {
		select {
		case <-server.ctx.Done():
			logger.Log.Print(2, "updateContainerStats stopped..")
			return
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(server.ctx, 5*time.Second)

			containers, err := server.service.ContainerList(ctx)
			if err != nil {
				logger.Log.Error("[BR] Service Container list error.. [%v]", err)
				cancel()
				continue
			}

			var wg sync.WaitGroup
			var ch_res chan docker.ContainerStats = make(chan docker.ContainerStats, len(containers))

			for _, c := range containers {
				c := c
				wg.Add(1)

				go func() {
					defer wg.Done()
					rst, err := server.service.ContainerStats(ctx, c.ID, false)
					if err != nil {
						logger.Log.Error("[upd] get Containerstats error [%s] [%v]", c.ID, err)
						return
					}

					if rst == nil {
						logger.Log.Error("[upd] get Containerstats rst is nil..(%s) ", c.ID)
						return
					}

					rst.ID = c.ID
					rst.Name = c.Name

					select {
					case ch_res <- *rst:
					case <-ctx.Done():
					}
				}()
			}

			go func() {
				wg.Wait()
				close(ch_res)
			}()

			select {
			case <-ctx.Done():
				logger.Log.Warn("[upd] ContainerStats timeout (5s)")
			}

			cancel()
			// close(ch_res)

			var res []docker.ContainerStats
			for r := range ch_res {
				res = append(res, r)
			}

			for _, r := range res {
				logger.Log.Print(2, "res : %v", r)
			}

			encoding, err := json.Marshal(res)
			server.hub.Broadcast(encoding)
		}
	}
}

// func (server *Server) updateContainerStats() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

// 	defer cancel()

// 	for {
// 		containers, err := server.service.ContainerList(ctx)
// 		if err != nil {
// 			logger.Log.Error("[BR] Service Container list error.. [%v]", err)
// 			return
// 		}

// 		var res []docker.ContainerStats = []docker.ContainerStats{}
// 		var wg sync.WaitGroup
// 		var ch_res chan docker.ContainerStats = make(chan docker.ContainerStats, len(containers))

// 		for _, c := range containers {
// 			c := c // loop 변수 캡처 방지
// 			wg.Add(1)
// 			logger.Log.Print(2, "container stats... %s", c.ID)

// 			go func() {
// 				defer wg.Done()
// 				rst, err := server.service.ContainerStats(ctx, c.ID, false)
// 				if err != nil {
// 					logger.Log.Error("[upd]Service statContainer error.. [%s] [%v]", c.ID, err)
// 				}
// 				select {
// 				case ch_res <- *rst:
// 				case <-ctx.Done():
// 				}
// 			}()
// 		}

// 		// wg타임아웃 5초
// 		done := make(chan struct{})
// 		go func() {
// 			wg.Wait()
// 			close(done)
// 		}()

// 		select {
// 		case <-done: // 정상종료
// 		case <-ctx.Done():
// 			logger.Log.Warn("[upd] ContainerStats timeout (5s)")
// 		}

// 		cancel()
// 		close(ch_res)

// 		for c := range ch_res {
// 			res = append(res, c)
// 		}

// 		for _, r := range res {
// 			logger.Log.Print(2, "res : %v", r)
// 		}

// 		time.Sleep(time.Second * 1)
// 	}
// }

func (server *Server) Start() error {
	logger.Log.Print(2, "Gin server start.")

	go server.hub.Run()
	go server.updateContainerStats()

	if err := server.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Error("listen error. %v", err)
		return err
	}

	return nil
}

func (server *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer server.wg.Done()
	server.svr_cancel()

	if err := server.srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Server Shutdown:", err)
		return err
	}
	return nil
}
