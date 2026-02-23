package api

import (
	"api-gateway/internal/config"
	"api-gateway/internal/container"
	"api-gateway/internal/db"
	"api-gateway/internal/logger"
	"api-gateway/internal/memory"
	"api-gateway/internal/service"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

// var serviceMap = map[string]string{
// 	"/auth":   "http://localhost:9082",
// 	"/docker": "http://localhost:9083",
// }

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

func newRESTProxy(target string) *httputil.ReverseProxy {
	url, _ := url.Parse(target)
	return httputil.NewSingleHostReverseProxy(url)
}

// func newSSEProxy(target string) *httputil.ReverseProxy {
// 	url, _ := url.Parse(target)

// 	proxy := httputil.NewSingleHostReverseProxy(url)

// 	proxy.Transport = &http.Transport{
// 		ForceAttemptHTTP2:     false,
// 		ResponseHeaderTimeout: 0,
// 	}

// 	proxy.ModifyResponse = func(resp *http.Response) error {
// 		resp.Header.Set("Content-Type", "text/event-stream")
// 		resp.Header.Set("Cache-Control", "no-cache")
// 		resp.Header.Set("Connection", "keep-alive")
// 		return nil
// 	}

// 	return proxy
// }

func (server *Server) newSSEProxy(target string) *httputil.ReverseProxy {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.Transport = &http.Transport{
		ForceAttemptHTTP2:     false, // 🔥 이거 없으면 거의 100% 끊김
		DisableKeepAlives:     false,
		ResponseHeaderTimeout: 0,
	}

	proxy.FlushInterval = -1

	proxy.ModifyResponse = func(resp *http.Response) error {
		resp.Header.Set("Content-Type", "text/event-stream")
		resp.Header.Set("Cache-Control", "no-cache")
		resp.Header.Set("Connection", "keep-alive")
		resp.Header.Set("Access-Control-Allow-Origin", "*")
		return nil
	}

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Log.Print(2, "SSE proxy error: %v", err)
	}

	return proxy
}

// func (server *Server) newReverseProxy(target string, w http.ResponseWriter) *httputil.ReverseProxy {
func (server *Server) newReverseProxy(target string, c *gin.Context) *httputil.ReverseProxy {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ModifyResponse = func(resp *http.Response) error {
		if resp.Request.URL.Path == "/login" && resp.StatusCode == http.StatusOK {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			resp.Body.Close()

			var loginResp loginUserResponse
			if err := json.Unmarshal(body, &loginResp); err != nil {
				logger.Log.Error("[newReverseProxy] /login parse error: %v", err)
			} else {
				logger.Log.Print(2, "[newReverseProxy] /login caught")
				logger.Log.Print(2, "\t ss:%s", loginResp.SessionID)
				logger.Log.Print(2, "\t at:%s", loginResp.AcessToken)
				logger.Log.Print(2, "\t ate:%s", loginResp.AccessTokenExpiresAt)
				logger.Log.Print(2, "\t rt:%s", loginResp.RefreshToken)
				logger.Log.Print(2, "\t rte:%s", loginResp.RefreshTokenExpiresAt)

				server.SetCookie(c, "refresh_token", loginResp.RefreshToken, loginResp.RefreshTokenExpiresAt)
			}

			// Body 복원 (클라이언트에 그대로 전달)
			resp.Body = io.NopCloser(bytes.NewBuffer(body))
			resp.ContentLength = int64(len(body))
		}
		return nil
	}

	return proxy
}

func (server *Server) setupRouter() {
	logger.Log.Print(2, "auth url : %v", server.config.AUTH_SERVICE_URL)
	logger.Log.Print(2, "docker url : %v", server.config.DOCKER_SERVICE_URL)

	// router := gin.Default()
	router := gin.New()
	addresses := strings.Split(server.config.AllowOrigins, ",")
	router.Use(corsMiddleware(addresses))
	router.Use(authMiddleware(server.tokenMaker))

	// gin.SetMode(gin.DebugMode)
	// fmt.Printf("%v, \n", server.config.AllowOrigins)

	router.GET("/heartbeat", server.heartbeat)
	router.GET("/terminate", server.terminate)

	// prefix 단위 라우팅
	router.Any("/auth/*proxyPath", func(c *gin.Context) {
		addr := server.config.AUTH_SERVICE_URL
		// proxy := server.newReverseProxy(addr, c.Writer)
		proxy := server.newReverseProxy(addr, c)
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/auth")
		logger.Log.Print(2, "auth path : %s", c.Request.URL.Path)
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	router.Any("/docker/*proxyPath", func(c *gin.Context) {
		// addr := server.config.DOCKER_SERVICE_URL
		// proxy := newReverseProxy(serviceMap["/docker"])

		logger.Log.Print(2, "docker url :%v ", c.Request.URL)

		// proxy := newReverseProxy(addr)
		proxy := newRESTProxy(server.config.DOCKER_SERVICE_URL)
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/docker")
		logger.Log.Print(2, "docker path : %s", c.Request.URL.Path)
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	// SSE 전용
	// router.GET("/docker-sse/events", func(c *gin.Context) {
	// 	logger.Log.Print(2, "docker-sse url :%v ", c.Request.URL)

	// 	proxy := newSSEProxy(server.config.DOCKER_SERVICE_URL)
	// 	proxy.FlushInterval = -1
	// 	c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/docker-sse")
	// 	logger.Log.Print(2, "dockersse path : %s", c.Request.URL.Path)
	// 	proxy.ServeHTTP(c.Writer, c.Request)
	// })

	// SSE 전용
	router.GET("/docker-sse/events", func(c *gin.Context) {
		logger.Log.Print(2, "docker-sse url :%v ", c.Request.URL)

		// WriteTimeout 비활성화
		rc := http.NewResponseController(c.Writer)
		rc.SetWriteDeadline(time.Time{}) // deadline 제거

		proxy := server.newSSEProxy(server.config.DOCKER_SERVICE_URL)
		proxy.FlushInterval = -1
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/docker-sse")
		logger.Log.Print(2, "dockersse path : %s", c.Request.URL.Path)
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
