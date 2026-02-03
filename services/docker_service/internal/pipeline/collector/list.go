package collector

import (
	"context"
	"sync"
	"time"

	"docker_service/internal/docker"
	"docker_service/internal/logger"
	"docker_service/internal/pipeline"
)

// ListCollector Container List 수집기
type ListCollector struct {
	client   *docker.Client
	config   Config
	buffer   *RingBuffer
	stopCh   chan struct{}
	stopOnce sync.Once
	wg       sync.WaitGroup
}

// NewListCollector ListCollector 생성
func NewListCollector(client *docker.Client, cfg Config) *ListCollector {
	return &ListCollector{
		client: client,
		config: cfg,
		buffer: NewRingBuffer(cfg.BufferSize),
		stopCh: make(chan struct{}),
	}
}

func (c *ListCollector) Name() string {
	return "list-collector"
}

func (c *ListCollector) Start(ctx context.Context) (<-chan pipeline.Message, error) {
	c.wg.Add(1)
	go c.run(ctx)
	return c.buffer.Channel(), nil
}

func (c *ListCollector) Stop() error {
	c.stopOnce.Do(func() {
		close(c.stopCh)
	})
	c.wg.Wait()
	c.buffer.Close()
	return nil
}

func (c *ListCollector) run(ctx context.Context) {
	defer c.wg.Done()

	ticker := time.NewTicker(time.Duration(c.config.IntervalSec) * time.Second)
	defer ticker.Stop()

	// 시작 시 즉시 한 번 수집
	c.collect(ctx)

	for {
		select {
		case <-ticker.C:
			c.collect(ctx)
		case <-c.stopCh:
			logger.Log.Print(2, "[ListCollector] stopped")
			return
		case <-ctx.Done():
			logger.Log.Print(2, "[ListCollector] context cancelled")
			return
		}
	}
}

func (c *ListCollector) collect(ctx context.Context) {
	containers, err := c.client.ListContainers(ctx)
	if err != nil {
		logger.Log.Error("[ListCollector] failed to list containers: %v", err)
		return
	}

	// docker.Container -> pipeline.ContainerInfo 변환
	infos := make([]pipeline.ContainerInfo, 0, len(containers))
	for _, ct := range containers {
		infos = append(infos, pipeline.ContainerInfo{
			ID:     ct.ID,
			Name:   ct.Name,
			Image:  ct.Image,
			State:  ct.State,
			Status: ct.Status,
		})
	}

	msg := pipeline.Message{
		Type:      pipeline.DataTypeList,
		Host:      c.config.Host,
		Timestamp: time.Now(),
		Data: pipeline.ContainerListData{
			Containers: infos,
		},
	}

	c.buffer.Send(msg)
	logger.Log.Print(2, "[ListCollector] collected %d containers from %s", len(containers), c.config.Host)
}

// CollectOnce 단발성 수집 (즉시 수집이 필요할 때)
func (c *ListCollector) CollectOnce(ctx context.Context) (*pipeline.Message, error) {
	containers, err := c.client.ListContainers(ctx)
	if err != nil {
		return nil, err
	}

	infos := make([]pipeline.ContainerInfo, 0, len(containers))
	for _, ct := range containers {
		infos = append(infos, pipeline.ContainerInfo{
			ID:     ct.ID,
			Name:   ct.Name,
			Image:  ct.Image,
			State:  ct.State,
			Status: ct.Status,
		})
	}

	msg := &pipeline.Message{
		Type:      pipeline.DataTypeList,
		Host:      c.config.Host,
		Timestamp: time.Now(),
		Data: pipeline.ContainerListData{
			Containers: infos,
		},
	}

	return msg, nil
}
