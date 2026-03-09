package event

import (
	"context"
	"docker_service/internal/config"
	"docker_service/internal/container"
	"docker_service/internal/event2"
	"docker_service/internal/logger"
	"sync"
)

type Server struct {
	ctx      context.Context
	cancel   context.CancelFunc
	wg       *sync.WaitGroup
	config   *config.Config
	eventMgr *event2.EventManager
}

func NewServer(wg *sync.WaitGroup, ct *container.Container, eventMgr *event2.EventManager) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		ctx:      ctx,
		cancel:   cancel,
		wg:       wg,
		config:   ct.Config,
		eventMgr: eventMgr,
	}, nil
}

func (s *Server) Start() error {
	// 1. EventManager 시작
	s.eventMgr.Start(s.ctx)

	// 2. 초기 호스트들 watch 시작 (설정된 모든 호스트)
	hosts, _ := s.config.GetDockerHosts()

	for _, host := range hosts {
		if err := s.eventMgr.WatchHost(host.Name); err != nil {
			logger.Log.Error("Failed to watch host %s: %v", host, err)
		}
	}
	return nil
}

func (s *Server) Shutdown() error {
	logger.Log.Print(3, "[EventManager] shutdown ...")
	defer s.wg.Done()

	s.cancel()
	// EventManager 종료 (graceful)
	s.eventMgr.Stop()

	logger.Log.Print(3, "[EventManager] shutdown complete")
	return nil
}
