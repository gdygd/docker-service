package api

import (
	"api-gateway/internal/config"
	"api-gateway/internal/container"
	"api-gateway/internal/db"
	"api-gateway/internal/logger"
	"api-gateway/internal/memory"
	"api-gateway/internal/service"
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	apiserv "api-gateway/internal/service/api"

	"github.com/gdygd/goglib/token"

	"github.com/gin-gonic/gin"
)

const (
	R_TIME_OUT = 5 * time.Second
	W_TIME_OUT = 5 * time.Second
)

var serviceMap = map[string]string{
	"/auth":   "http://localhost:9082",
	"/docker": "http://localhost:9083",
}

// Server serves HTTP requests for our banking service.
type Server struct {
	wg           *sync.WaitGroup
	srv          *http.Server
	config       *config.Config
	tokenMaker   token.Maker
	router       *gin.Engine
	service      service.ServiceInterface
	dbHnd        db.DbHandler
	objdb        *memory.RedisDb
	ch_terminate chan bool
}

func NewServer(wg *sync.WaitGroup, ct *container.Container, ch_terminate chan bool) (*Server, error) {
	// init service
	apiservice := apiserv.NewApiService(ct.DbHnd, ct.ObjDb)
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
		objdb:        ct.ObjDb,
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

func newReverseProxy(target string) *httputil.ReverseProxy {
	url, _ := url.Parse(target)
	return httputil.NewSingleHostReverseProxy(url)
}

func (server *Server) setupRouter() {
	router := gin.Default()
	// gin.SetMode(gin.DebugMode)
	// fmt.Printf("%v, \n", server.config.AllowOrigins)

	addresses := strings.Split(server.config.AllowOrigins, ",")

	router.GET("/heartbeat", server.heartbeat)
	router.GET("/terminate", server.terminate)

	router.Use(corsMiddleware(addresses))
	// router.Use(authMiddleware(server.tokenMaker))

	// prefix 단위 라우팅
	router.Any("/auth/*proxyPath", func(c *gin.Context) {
		proxy := newReverseProxy(serviceMap["/auth"])
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/auth")
		logger.Log.Print(2, "auth path : %s", c.Request.URL.Path)
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	router.Any("/docker/*proxyPath", func(c *gin.Context) {
		proxy := newReverseProxy(serviceMap["/docker"])
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/docker")
		logger.Log.Print(2, "docker path : %s", c.Request.URL.Path)
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	server.router = router
}

func (server *Server) Start() error {
	logger.Log.Print(2, "Gin server start.")

	if err := server.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Log.Error("listen error. %v", err)
		return err
	}
	// if err := server.srv.ListenAndServeTLS("./tls/server.crt", "./tls/server.key"); err != nil && err != http.ErrServerClosed {
	// 	logger.Log.Error("listen error. %v", err)
	// 	return err
	// }

	return nil
}

func (server *Server) Shutdown() error {
	logger.Log.Print(2, "ShutDown..#1")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer server.wg.Done()
	if err := server.srv.Shutdown(ctx); err != nil {
		logger.Log.Error("Server Shutdown:", err)
		return err
	}
	logger.Log.Print(2, "ShutDown..#2")
	return nil
}
