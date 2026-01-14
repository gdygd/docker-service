package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"docker_service/internal/config"
	"docker_service/internal/container"
	"docker_service/internal/db"
	"docker_service/internal/logger"
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
	wg           *sync.WaitGroup
	srv          *http.Server
	config       *config.Config
	tokenMaker   token.Maker
	router       *gin.Engine
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

	server := &Server{
		wg:           wg,
		config:       ct.Config,
		tokenMaker:   tokenMaker,
		service:      apiservice,
		dbHnd:        ct.DbHnd,
		ch_terminate: ch_terminate,
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

func (server *Server) Start() error {
	logger.Log.Print(2, "Gin server start.")

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
	if err := server.srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Server Shutdown:", err)
		return err
	}
	return nil
}
