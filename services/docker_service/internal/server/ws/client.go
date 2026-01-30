package ws

import (
	"fmt"
	"time"

	"docker_service/internal/logger"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = 50 * time.Second
)

type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}

func (c *Client) WsRead() {
	// defer func() {
	// 	c.Hub.Unregister <- c
	// }()
	defer func() {
		select {
		case c.Hub.Unregister <- c:
		case <-c.Hub.Ctx.Done():
			// hub 종료 중 → 보내지 않기위해
		}
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		logger.Log.Print(2, "read#1 ")
		_, msg, err := c.Conn.ReadMessage()
		logger.Log.Print(2, "read#2 ")
		if err != nil {
			logger.Log.Print(2, "[client] read error: %v", err)
			return
		}
		logger.Log.Print(2, "[client] recv:%s", string(msg))

		res := fmt.Sprintf("%s hello", msg)
		c.Send <- []byte(res)
	}
}

func (c *Client) WsWrite() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		// c.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				logger.Log.Error("close message...(%v)", msg)
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				logger.Log.Error("write message error : %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		}
	}
}

func (c *Client) Close() {
	// close(c.Send)
	c.Conn.Close()
}
