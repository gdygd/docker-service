package api

import (
	"net/http"

	"docker_service/internal/logger"
	"docker_service/internal/server/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket 업그레이더
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (server *Server) wsHandler(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Log.Print(2, "[ws] ws upgrade error: %v", err)
		return
	}

	logger.Log.Print(2, "request client..")

	client := &ws.Client{
		Hub:  server.hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	server.hub.Register <- client

	go client.WsRead()
	go client.WsWrite()
}

// func (server *Server) wsHandler(c *gin.Context) {
// 	containerID := c.Param("id")

// 	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		conn.WriteJSON(gin.H{"error": err.Error()})
// 		return
// 	}
// 	defer conn.Close()

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	var ch_rst chan *docker.ContainerStats

// 	server.service.ContainerStatsStream(ctx, containerID, true, ch_rst)

// 	for {
// 		select {
// 		case stats := <-ch_rst:
// 			logger.Log.Print(2, "cpu: %f", stats.CPUPercent)

// 			if err := conn.WriteJSON(stats); err != nil {
// 				conn.WriteJSON(gin.H{"error": err.Error()})
// 				return
// 			}
// 		}
// 	}
// }
