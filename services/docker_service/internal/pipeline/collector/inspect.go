package collector

import (
	"context"
	"docker_service/internal/docker"
	"docker_service/internal/logger"
	"docker_service/internal/pipeline"
	"sync"
	"time"
)

// InspectCollector Container Inspect 수집기
type InspectCollector struct {
	client   *docker.Client
	config   Config
	buffer   *RingBuffer
	stopCh   chan struct{}
	stopOnce sync.Once
	wg       sync.WaitGroup // 루틴 종료 대기
}

// NewInspectCollector InspectCollector 생성
func NewInspectCollector(client *docker.Client, cfg Config) *InspectCollector {
	return &InspectCollector{
		client: client,
		config: cfg,
		buffer: NewRingBuffer(cfg.BufferSize),
		stopCh: make(chan struct{}),
	}
}

func (c *InspectCollector) Name() string {
	return "inspect-collector"
}

func (c *InspectCollector) Start(ctx context.Context) (<-chan pipeline.Message, error) {
	c.wg.Add(1)
	go c.run(ctx)
	return c.buffer.Channel(), nil
}

func (c *InspectCollector) Stop() error {
	c.stopOnce.Do(func() {
		close(c.stopCh)
	})
	c.wg.Wait()
	c.buffer.Close()
	return nil
}

func (c *InspectCollector) run(ctx context.Context) {
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
			logger.Log.Print(2, "[InspectCollector] stopped")
			return
		case <-ctx.Done():
			logger.Log.Print(2, "[InspectCollector] context cancelled")
			return
		}
	}
}

func (c *InspectCollector) collect(ctx context.Context) {
	// 먼저 컨테이너 목록 조회
	containers, err := c.client.ListContainers(ctx)
	if err != nil {
		logger.Log.Error("[InspectCollector] failed to list containers: %v", err)
		return
	}

	inspects := make([]pipeline.ContainerInspectInfo, 0, len(containers))

	// 각 컨테이너의 inspect 수집  (추후 동시 작업을 위한 worker 고려.)
	for _, ct := range containers {
		inspectResult, err := c.client.InspectContainer(ctx, ct.ID)
		if err != nil {
			logger.Log.Error("[InspectCollector] failed to inspect %s: %v", ct.ID, err)
			continue
		}

		// SDK 결과 -> docker.ContainerInspect -> pipeline.ContainerInspectInfo 변환
		dockerInspect := docker.ConvertInspectResult(inspectResult)
		info := convertInspectResult(dockerInspect)
		inspects = append(inspects, info)
	}

	msg := pipeline.Message{
		Type:      pipeline.DataTypeInspect,
		Host:      c.config.Host,
		Timestamp: time.Now(),
		Data: pipeline.ContainerInspectData{
			Inspects: inspects,
		},
	}

	c.buffer.Send(msg)
	logger.Log.Print(2, "[InspectCollector] collected %d inspects from %s", len(inspects), c.config.Host)
}

// convertInspectResult docker API 결과를 pipeline 타입으로 변환
func convertInspectResult(result docker.ContainerInspect) pipeline.ContainerInspectInfo {
	info := pipeline.ContainerInspectInfo{
		ID:           result.ID,
		Name:         result.Name,
		Image:        result.Image,
		Created:      result.Created,
		Platform:     result.Platform,
		RestartCount: result.RestartCount,
	}

	// State 변환
	if result.State != nil {
		info.State = &pipeline.ContainerStateInfo{
			Status:     result.State.Status,
			Running:    result.State.Running,
			Paused:     result.State.Paused,
			Restarting: result.State.Restarting,
			ExitCode:   result.State.ExitCode,
			StartedAt:  result.State.StartedAt,
			FinishedAt: result.State.FinishedAt,
		}
	}

	// Config 변환
	if result.Config != nil {
		info.Config = &pipeline.ContainerConfigInfo{
			Hostname:   result.Config.Hostname,
			User:       result.Config.User,
			Env:        result.Config.Env,
			Cmd:        result.Config.Cmd,
			Entrypoint: result.Config.Entrypoint,
			WorkingDir: result.Config.WorkingDir,
			Labels:     result.Config.Labels,
		}
	}

	// Network 변환
	if result.NetworkSettings != nil {
		ports := make(map[string][]pipeline.PortBindingInfo)
		for portKey, bindings := range result.NetworkSettings.Ports {
			portBindings := make([]pipeline.PortBindingInfo, 0, len(bindings))
			for _, b := range bindings {
				portBindings = append(portBindings, pipeline.PortBindingInfo{
					HostIP:   b.HostIP,
					HostPort: b.HostPort,
				})
			}
			ports[portKey] = portBindings
		}

		networks := make(map[string]pipeline.NetworkEndpoint)
		for netName, endpoint := range result.NetworkSettings.Networks {
			networks[netName] = pipeline.NetworkEndpoint{
				NetworkID:  endpoint.NetworkID,
				IPAddress:  endpoint.IPAddress,
				Gateway:    endpoint.Gateway,
				MacAddress: endpoint.MacAddress,
			}
		}

		info.Network = &pipeline.ContainerNetworkInfo{
			IPAddress:  result.NetworkSettings.IPAddress,
			Gateway:    result.NetworkSettings.Gateway,
			MacAddress: result.NetworkSettings.MacAddress,
			Ports:      ports,
			Networks:   networks,
		}
	}

	// Mounts 변환
	mounts := make([]pipeline.MountPointInfo, 0, len(result.Mounts))
	for _, m := range result.Mounts {
		mounts = append(mounts, pipeline.MountPointInfo{
			Type:        m.Type,
			Name:        m.Name,
			Source:      m.Source,
			Destination: m.Destination,
			Mode:        m.Mode,
			RW:          m.RW,
		})
	}
	info.Mounts = mounts

	return info
}
