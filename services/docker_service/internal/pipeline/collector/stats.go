package collector

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"docker_service/internal/docker"
	"docker_service/internal/logger"
	"docker_service/internal/pipeline"
)

// StatsCollector Container Stats 수집기
type StatsCollector struct {
	client   *docker.Client
	config   Config
	buffer   *RingBuffer
	stopCh   chan struct{}
	stopOnce sync.Once
	wg       sync.WaitGroup
}

// NewStatsCollector StatsCollector 생성
func NewStatsCollector(client *docker.Client, cfg Config) *StatsCollector {
	return &StatsCollector{
		client: client,
		config: cfg,
		buffer: NewRingBuffer(cfg.BufferSize),
		stopCh: make(chan struct{}),
	}
}

func (c *StatsCollector) Name() string {
	return "stats-collector"
}

func (c *StatsCollector) Start(ctx context.Context) (<-chan pipeline.Message, error) {
	c.wg.Add(1)
	go c.run(ctx)
	return c.buffer.Channel(), nil
}

func (c *StatsCollector) Stop() error {
	c.stopOnce.Do(func() {
		close(c.stopCh)
	})
	c.wg.Wait()
	c.buffer.Close()
	return nil
}

func (c *StatsCollector) run(ctx context.Context) {
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
			logger.Log.Print(2, "[StatsCollector] stopped")
			return
		case <-ctx.Done():
			logger.Log.Print(2, "[StatsCollector] context cancelled")
			return
		}
	}
}

func (c *StatsCollector) collect(ctx context.Context) {
	// 1. 컨테이너 목록 조회
	containers, err := c.client.ListContainers(ctx)
	if err != nil {
		logger.Log.Error("[StatsCollector] failed to list containers: %v", err)
		return
	}

	if len(containers) == 0 {
		logger.Log.Print(2, "[StatsCollector] no containers found")
		return
	}

	// 2. 결과 수집용 채널 및 WaitGroup
	var wg sync.WaitGroup
	resultCh := make(chan pipeline.ContainerStatsInfo, len(containers))

	// 3. 타임아웃 컨텍스트 (3초)
	childCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// 4. 각 컨테이너에 대해 goroutine으로 stats 수집
	for _, container := range containers {
		container := container // capture loop variable
		wg.Add(1)

		go func() {
			defer wg.Done()

			stats, err := c.getContainerStats(childCtx, container.ID)
			if err != nil {
				logger.Log.Error("[StatsCollector] failed to get stats for %s: %v", container.ID, err)
				return
			}

			if stats == nil {
				return
			}

			stats.ID = container.ID
			stats.Name = container.Name

			select {
			case resultCh <- *stats:
			case <-childCtx.Done():
			}
		}()
	}

	// 5. 완료 신호
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(resultCh)
		close(done)
	}()

	// 6. 완료 또는 타임아웃 대기
	select {
	case <-done:
		logger.Log.Print(2, "[StatsCollector] all goroutines completed")
	case <-childCtx.Done():
		logger.Log.Print(2, "[StatsCollector] timeout(3s)")
	}

	// 7. 결과 수집
	var statsInfos []pipeline.ContainerStatsInfo
	for stat := range resultCh {
		statsInfos = append(statsInfos, stat)
	}

	if len(statsInfos) == 0 {
		return
	}

	// 8. 메시지 생성 및 전송
	msg := pipeline.Message{
		Type:      pipeline.DataTypeStats,
		Host:      c.config.Host,
		Timestamp: time.Now(),
		Data: pipeline.ContainerStatsData{
			Stats: statsInfos,
		},
	}

	c.buffer.Send(msg)
	logger.Log.Print(2, "[StatsCollector] collected %d stats from %s", len(statsInfos), c.config.Host)
}

// getContainerStats 개별 컨테이너 Stats 조회
func (c *StatsCollector) getContainerStats(ctx context.Context, containerID string) (*pipeline.ContainerStatsInfo, error) {
	result, err := c.client.ContainerStats(ctx, containerID, true)
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	decoder := json.NewDecoder(result.Body)

	// 첫 프레임 (버림)
	var first docker.ContainerStatsRaw
	if err := decoder.Decode(&first); err != nil {
		return nil, err
	}

	// 두 번째 프레임 (이걸로 계산)
	var second docker.ContainerStatsRaw
	if err := decoder.Decode(&second); err != nil {
		return nil, err
	}

	return c.calculateStats(second), nil
}

// calculateStats raw 데이터를 ContainerStatsInfo로 변환
func (c *StatsCollector) calculateStats(raw docker.ContainerStatsRaw) *pipeline.ContainerStatsInfo {
	// CPU %
	cpuDelta := float64(
		raw.CPUStats.CPUUsage.TotalUsage -
			raw.PreCPUStats.CPUUsage.TotalUsage,
	)

	systemDelta := float64(
		raw.CPUStats.SystemCPUUsage -
			raw.PreCPUStats.SystemCPUUsage,
	)

	cpuPercent := 0.0
	if systemDelta > 0 && cpuDelta > 0 {
		cpuPercent = (cpuDelta / systemDelta) *
			float64(len(raw.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}

	// Memory
	memUsage := c.calculateMemoryUsage(raw)
	memLimit := raw.MemoryStats.Limit

	memPercent := 0.0
	if memLimit > 0 {
		memPercent = (float64(memUsage) / float64(memLimit)) * 100.0
	}

	// Network
	var rx, tx uint64
	for _, n := range raw.Networks {
		rx += n.RxBytes
		tx += n.TxBytes
	}

	return &pipeline.ContainerStatsInfo{
		CPUPercent:    cpuPercent,
		MemoryUsage:   memUsage,
		MemoryLimit:   memLimit,
		MemoryPercent: memPercent,
		NetworkRx:     rx,
		NetworkTx:     tx,
	}
}

// calculateMemoryUsage 캐시를 제외한 메모리 사용량 계산
func (c *StatsCollector) calculateMemoryUsage(raw docker.ContainerStatsRaw) uint64 {
	// cgroup v2: inactive_file 사용
	if raw.MemoryStats.Stats.InactiveFile > 0 {
		if raw.MemoryStats.Usage > raw.MemoryStats.Stats.InactiveFile {
			return raw.MemoryStats.Usage - raw.MemoryStats.Stats.InactiveFile
		}
	}
	// cgroup v1: cache 사용
	if raw.MemoryStats.Stats.Cache > 0 {
		if raw.MemoryStats.Usage > raw.MemoryStats.Stats.Cache {
			return raw.MemoryStats.Usage - raw.MemoryStats.Stats.Cache
		}
	}
	return raw.MemoryStats.Usage
}
