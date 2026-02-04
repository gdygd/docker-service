package pipe

import (
	"context"
	"docker_service/internal/docker"
	"docker_service/internal/logger"
	"docker_service/internal/pipeline"
	"docker_service/internal/pipeline/collector"
	"sync"
)

// Server Pipeline 데이터 수집 서버
type Server struct {
	ctx     context.Context
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
	manager *collector.Manager
	config  Config
}

// Config Pipeline 서버 설정
type Config struct {
	IntervalSec int // 수집 주기 (초)
	BufferSize  int // 채널 버퍼 크기
}

// DefaultConfig 기본 설정
func DefaultConfig() Config {
	return Config{
		IntervalSec: 30,
		BufferSize:  50,
	}
}

// NewServer Pipeline 서버 생성
func NewServer(wg *sync.WaitGroup, dockerMng *docker.DockerClientManager, cfg Config) (*Server, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Manager 생성
	manager := collector.NewManager(dockerMng, 100)

	// Collector 설정
	collectorCfg := collector.Config{
		IntervalSec: cfg.IntervalSec,
		BufferSize:  cfg.BufferSize,
	}

	// 모든 호스트에 Collector 등록
	if err := manager.RegisterAllHosts(
		[]collector.CollectorType{collector.TypeList, collector.TypeInspect},
		collectorCfg,
	); err != nil {
		logger.Log.Error("[PipeServer] collector registration fail: %v", err)
	}

	return &Server{
		ctx:     ctx,
		cancel:  cancel,
		wg:      wg,
		manager: manager,
		config:  cfg,
	}, nil
}

// Start Pipeline 서버 시작
func (s *Server) Start() error {
	outCh, err := s.manager.Start(s.ctx)
	if err != nil {
		logger.Log.Error("[PipeServer] start fail: %v", err)
		return err
	}

	logger.Log.Print(3, "[PipeServer] started, collectors: %d", s.manager.GetCollectorCount())

	// 메시지 처리 루프
	s.processMessages(outCh)

	return nil
}

// Shutdown Pipeline 서버 종료
func (s *Server) Shutdown() error {
	logger.Log.Print(3, "[PipeServer] shutting down...")
	defer s.wg.Done()

	s.cancel()
	s.manager.Stop()

	logger.Log.Print(3, "[PipeServer] shutdown complete")
	return nil
}

// GetCollectorCount 등록된 Collector 수 반환
func (s *Server) GetCollectorCount() int {
	return s.manager.GetCollectorCount()
}

// processMessages 수집된 메시지 처리
func (s *Server) processMessages(outCh <-chan pipeline.Message) {
	for msg := range outCh {
		s.handleMessage(msg)
	}
}

// handleMessage 개별 메시지 처리
func (s *Server) handleMessage(msg pipeline.Message) {
	// TODO: gRPC Sender로 전송
	logger.Log.Print(2, "[PipeServer] type=%s host=%s timestamp=%v",
		msg.Type, msg.Host, msg.Timestamp)

	switch msg.Type {
	case pipeline.DataTypeList:
		s.handleListMessage(msg)
	case pipeline.DataTypeInspect:
		s.handleInspectMessage(msg)
	}
}

// handleListMessage Container List 메시지 처리
func (s *Server) handleListMessage(msg pipeline.Message) {
	containers := msg.Data.(pipeline.ContainerListData)
	logger.Log.Print(2, "[PipeServer] Container List (%d )>> ", len(containers.Containers))
	for _, c := range containers.Containers {
		logger.Log.Print(1, "\t ID:%s, Name:%s, Image:%s, State:%s, Status:%s",
			c.ID, c.Name, c.Image, c.State, c.Status)
	}
}

// handleInspectMessage Container Inspect 메시지 처리
func (s *Server) handleInspectMessage(msg pipeline.Message) {
	inspects := msg.Data.(pipeline.ContainerInspectData)
	logger.Log.Print(2, "[PipeServer] Inspect List (%d )>> ", len(inspects.Inspects))
	for _, ins := range inspects.Inspects {
		logger.Log.Print(1, "#[basic] Id:%s, Name:%s, Created:%s Platform:%s",
			ins.ID, ins.Name, ins.Created, ins.Platform)

		if ins.State != nil {
			logger.Log.Print(1, "\t [state] status:%s, running:%v, exitcode:%d, startedat:%s",
				ins.State.Status, ins.State.Running, ins.State.ExitCode, ins.State.StartedAt)
		}

		if ins.Config != nil {
			logger.Log.Print(1, "\t [config] host:%s, user:%v, env:%v, cmd:%v",
				ins.Config.Hostname, ins.Config.User, ins.Config.Env, ins.Config.Cmd)
		}

		if ins.Network != nil {
			logger.Log.Print(1, "\t [network] ip:%s, gw:%v, mac:%v, port:%v",
				ins.Network.IPAddress, ins.Network.Gateway, ins.Network.MacAddress, ins.Network.Ports)
		}

		for _, m := range ins.Mounts {
			logger.Log.Print(1, "\t [mount] type:%s, name:%v, src:%v, dst:%v mode:%s, rw:%v",
				m.Type, m.Name, m.Source, m.Destination, m.Mode, m.RW)
		}
	}
}
