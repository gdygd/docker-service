package ws

import (
	"context"
	"sync/atomic"

	"docker_service/internal/logger"
)

type Hub struct {
	clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	broadcast  chan []byte
	Ctx        context.Context
	closing    atomic.Bool
}

func NewHub(ctx context.Context) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		broadcast:  make(chan []byte), // 전체 클라이언트 broadcase
		Ctx:        ctx,
	}
}

func (h *Hub) Broadcast(msg []byte) {
	logger.Log.Print(2, "broadcast.. %d", len(msg))
	h.broadcast <- msg
}

func (h *Hub) Run() {
	logger.Log.Print(2, "hub run...")
	for {
		select {
		case <-h.Ctx.Done():
			logger.Log.Print(2, "[hub] shutdown")
			h.closing.Store(true)

			for c := range h.clients {
				close(c.Send)
				c.Close()
			}
			return

		case c := <-h.Register:
			h.clients[c] = true
			logger.Log.Print(2, "[hub] client connected:", len(h.clients))

		case c := <-h.Unregister:
			if h.closing.Load() {
				// shutdown중이면 무시
				continue
			}

			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.Send)
				c.Close()
				logger.Log.Print(2, "[hub] client disconnected:", len(h.clients))
			}

		case msg := <-h.broadcast:
			var deadClients []*Client

			for c := range h.clients {
				select {
				case c.Send <- msg:
				default:

					// delete(h.clients, c)
					// close(c.Send)
					// c.Close()

					// send실패시 나중에 한번에 정리
					deadClients = append(deadClients, c)
				}
			}

			// deadclients 정리
			for _, c := range deadClients {
				if _, ok := h.clients[c]; ok {
					delete(h.clients, c)
					close(c.Send)
					c.Close()
				}
			}
		}
	}
}
