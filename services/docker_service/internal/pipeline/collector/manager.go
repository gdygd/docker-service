package collector

import (
	"context"
	"docker_service/internal/docker"
	"docker_service/internal/logger"
	"docker_service/internal/pipeline"
	"sync"
)

// CollectorType 수집기 타입
type CollectorType string

const (
	TypeList    CollectorType = "list"
	TypeInspect CollectorType = "inspect"
)

// Manager 멀티 호스트 Collector 관리자
type Manager struct {
	dockerMng  *docker.DockerClientManager
	collectors map[string][]Collector // key: host name
	outCh      chan pipeline.Message
	mu         sync.RWMutex
	wg         sync.WaitGroup
}

// NewManager Collector Manager 생성
func NewManager(dockerMng *docker.DockerClientManager, bufferSize int) *Manager {
	return &Manager{
		dockerMng:  dockerMng,
		collectors: make(map[string][]Collector),
		outCh:      make(chan pipeline.Message, bufferSize),
	}
}

// RegisterCollectors 특정 호스트에 대한 수집기 등록
func (m *Manager) RegisterCollectors(hostName string, types []CollectorType, cfg Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, err := m.dockerMng.Get(hostName)
	if err != nil {
		return err
	}

	cfg.Host = hostName
	var collectors []Collector

	for _, t := range types {
		var c Collector
		switch t {
		case TypeList:
			c = NewListCollector(client, cfg)
		case TypeInspect:
			c = NewInspectCollector(client, cfg)
		}
		if c != nil {
			collectors = append(collectors, c)
		}
	}

	m.collectors[hostName] = collectors
	return nil
}

// RegisterAllHosts 모든 등록된 Docker 호스트에 수집기 등록
func (m *Manager) RegisterAllHosts(types []CollectorType, cfg Config) error {
	hostNames := m.dockerMng.GetHostNames()

	for _, hostName := range hostNames {
		if err := m.RegisterCollectors(hostName, types, cfg); err != nil {
			logger.Log.Error("[CollectorManager] failed to register collectors for %s: %v", hostName, err)
			continue
		}
		logger.Log.Print(2, "[CollectorManager] registered collectors for host: %s", hostName)
	}

	return nil
}

// Start 모든 수집기 시작
func (m *Manager) Start(ctx context.Context) (<-chan pipeline.Message, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for hostName, collectors := range m.collectors {
		for _, c := range collectors {
			ch, err := c.Start(ctx)
			if err != nil {
				logger.Log.Error("[CollectorManager] failed to start %s for %s: %v",
					c.Name(), hostName, err)
				continue
			}

			// 각 수집기의 출력을 통합 채널로 전달
			m.wg.Add(1)
			go func(ch <-chan pipeline.Message, name string) {
				defer m.wg.Done()
				for msg := range ch {
					select {
					case m.outCh <- msg:
					case <-ctx.Done():
						return
					}
				}
			}(ch, c.Name())

			logger.Log.Print(2, "[CollectorManager] started %s for host: %s", c.Name(), hostName)
		}
	}

	return m.outCh, nil
}

// Stop 모든 수집기 중지
func (m *Manager) Stop() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for hostName, collectors := range m.collectors {
		for _, c := range collectors {
			if err := c.Stop(); err != nil {
				logger.Log.Error("[CollectorManager] failed to stop %s for %s: %v",
					c.Name(), hostName, err)
			}
		}
	}

	m.wg.Wait()
	close(m.outCh)
	return nil
}

// GetCollectorCount 등록된 수집기 수 반환
func (m *Manager) GetCollectorCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, collectors := range m.collectors {
		count += len(collectors)
	}
	return count
}
