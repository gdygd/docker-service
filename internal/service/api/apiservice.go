package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"docker_service/internal/db"
	"docker_service/internal/docker"
	"docker_service/internal/logger"
	"docker_service/internal/service"
)

const (
	_          = iota
	KiB uint64 = 1 << (10 * iota)
	MiB
	GiB
)

type ApiService struct {
	dbHnd  db.DbHandler
	docker *docker.Client
	docMng *docker.DockerClientManager
}

func NewApiService(dbHnd db.DbHandler, docker *docker.Client, dockerMng *docker.DockerClientManager) service.ServiceInterface {
	return &ApiService{
		dbHnd:  dbHnd,
		docker: docker,    // none tls client	(only local host)
		docMng: dockerMng, // tls client	(for remote and local host)
	}
}

func (s *ApiService) Test() {
	fmt.Printf("test service")
}

func (s *ApiService) ContainerList(ctx context.Context) ([]docker.Container, error) {
	return s.docker.ListContainers(ctx)
}

func (s *ApiService) ContainerList2(ctx context.Context, host string) ([]docker.Container, error) {
	client, err := s.docMng.Get(host)
	if err != nil {
		logger.Log.Error("[ContainerList2] Get host client error..(%v)", err)
	}

	return client.ListContainers(ctx)
}

func (s *ApiService) InspectContainer(ctx context.Context, containerID string) (docker.ContainerInspect, error) {
	res, err := s.docker.InspectContainer(ctx, containerID)
	if err != nil {
		logger.Log.Error("inspect container error.. .%v", err)
		return docker.ContainerInspect{}, err
	}

	logger.Log.Print(2, "ID: %s", res.Container.ID)
	logger.Log.Print(2, "Image : %s", res.Container.Image)
	logger.Log.Print(2, "Name: %s", res.Container.Name)

	return docker.ContainerInspect{
		ID:    res.Container.ID,
		Image: res.Container.Image,
		Name:  res.Container.Name,
	}, nil
}

func (s *ApiService) InspectContainer2(ctx context.Context, containerID, host string) (docker.ContainerInspect, error) {
	client, err := s.docMng.Get(host)
	if err != nil {
		logger.Log.Error("[InspectContainer2] Get host client error..(%v)", err)
	}

	res, err := client.InspectContainer(ctx, containerID)
	if err != nil {
		logger.Log.Error("InspectContainer2 error.. .%v", err)
		return docker.ContainerInspect{}, err
	}

	logger.Log.Print(2, "ID: %s", res.Container.ID)
	logger.Log.Print(2, "Image : %s", res.Container.Image)
	logger.Log.Print(2, "Name: %s", res.Container.Name)

	return docker.ContainerInspect{
		ID:    res.Container.ID,
		Image: res.Container.Image,
		Name:  res.Container.Name,
	}, nil
}

func (s *ApiService) StartContainer(ctx context.Context, id string) error {
	rst, err := s.docker.StartContainer(ctx, id)
	if err != nil {
		logger.Log.Error("StartContainer err .. %v", err)
		return err
	}

	logger.Log.Print(2, "StartContainer rst : %v", rst)
	return nil
}

func (s *ApiService) StopContainer(ctx context.Context, id string) error {
	rst, err := s.docker.StopContainer(ctx, id)
	if err != nil {
		logger.Log.Error("StopContainer err .. %v", err)
		return err
	}

	logger.Log.Print(2, "StopContainer rst : %v", rst)
	return nil
}

// func (s *ApiService) ContainerStats(ctx context.Context, id string, stream bool) (*docker.ContainerStats, error) {
// 	result, err := s.docker.ContainerStats(ctx, id, stream)
// 	if err != nil {
// 		logger.Log.Error("Get ContainerStats error.. (%v)", err)
// 		return nil, err
// 	}
// 	defer result.Body.Close()

// 	var raw docker.ContainerStatsRaw
// 	if err := json.NewDecoder(result.Body).Decode(&raw); err != nil {
// 		logger.Log.Error("ContainerStats raw data decode error.. (%v)", err)
// 		return nil, err
// 	}

// 	return calculateStats(raw), nil
// }

func (s *ApiService) ContainerStats(ctx context.Context, id string, stream bool) (*docker.ContainerStats, error) {
	result, err := s.docker.ContainerStats(ctx, id, true)
	if err != nil {
		logger.Log.Error("Get ContainerStats error.. (%v)", err)
		return nil, err
	}
	defer result.Body.Close()

	decoder := json.NewDecoder(result.Body)

	// 첫 프레임 (버림)
	var first docker.ContainerStatsRaw
	if err := decoder.Decode(&first); err != nil {
		logger.Log.Error("ContainerStats raw data decode error #1.. (%v)", err)
		return nil, err
	}

	// 두 번째 프레임 (이걸로 계산)
	var second docker.ContainerStatsRaw
	if err := decoder.Decode(&second); err != nil {
		logger.Log.Error("ContainerStats raw data decode error #2.. (%v)", err)
		return nil, err
	}

	stats := calculateStats(second)
	return stats, nil
}

func (s *ApiService) ContainerStatsStream(ctx context.Context, id string, stream bool, ch_rst chan *docker.ContainerStats) error {
	result, err := s.docker.ContainerStats(ctx, id, stream)
	if err != nil {
		logger.Log.Error("Get ContainerStatsStream error.. (%v)", err)
		return err
	}
	defer result.Body.Close()
	decoder := json.NewDecoder(result.Body)

	for {
		select {
		case <-ctx.Done():
			logger.Log.Print(2, "Stop ContainerStatsStream..")
			return nil
		default:
			var statsRaw docker.ContainerStatsRaw
			if err := decoder.Decode(&statsRaw); err != nil {
				logger.Log.Error("ContainerStats raw data decode error #1.. (%v)", err)
				return err
			}

			stats := calculateStats(statsRaw)
			ch_rst <- stats

			time.Sleep(time.Second * 1)
		}
	}
}

func calculateStats(raw docker.ContainerStatsRaw) *docker.ContainerStats {
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

	// Memory %
	// memUsage := raw.MemoryStats.Usage	// 캐시 값을 제외하지 않은, 전체 메모리 usage
	memUsage := calculateMemoryUsage(raw) // 캐시 값을 제외한 메모리 usage
	memLimit := raw.MemoryStats.Limit

	memPercent := 0.0
	if memLimit > 0 {
		memPercent = (float64(memUsage) / float64(memLimit)) * 100.0
	}

	// 포맷팅 적용 Memory
	// usageVal, usageUnit := formatBytes(raw.MemoryStats.Usage)
	usageVal, usageUnit := formatBytes(memUsage)
	limitVal, limitUnit := formatBytes(raw.MemoryStats.Limit)

	// Network
	var rx, tx uint64
	for _, n := range raw.Networks {
		rx += n.RxBytes
		tx += n.TxBytes
	}

	return &docker.ContainerStats{
		CPUPercent:      cpuPercent,
		MemoryUsage:     memUsage,
		MemoryLimit:     memLimit,
		MemoryUsageVal:  round(usageVal, 2),
		MemoryUsageUnit: usageUnit,
		MemoryLimitVal:  round(limitVal, 2),
		MemoryLimitUnit: limitUnit,
		MemoryPercent:   memPercent,
		NetworkRx:       rx,
		NetworkTx:       tx,
	}
}

// func calculateMemoryUsage(raw docker.ContainerStatsRaw) uint64 {
// 	usage := raw.MemoryStats.Usage

// 	// cgroup v1
// 	if raw.MemoryStats.Stats.Cache > 0 {
// 		return usage - raw.MemoryStats.Stats.Cache
// 	}

// 	// cgroup v2
// 	if raw.MemoryStats.Stats.InactiveFile > 0 {
// 		return usage - raw.MemoryStats.Stats.InactiveFile
// 	}

// 	return usage
// }

func calculateMemoryUsage(raw docker.ContainerStatsRaw) uint64 {
	usage := raw.MemoryStats.Usage

	cgroup := detectCgroupVersion()

	if cgroup == 2 {
		return usage - raw.MemoryStats.Stats.Cache
	} else {
		return usage - raw.MemoryStats.Stats.InactiveFile
	}
}

func detectCgroupVersion() int {
	if _, err := os.Stat("/sys/fs/cgroup/cgroup.controllers"); err == nil {
		return 2 // cgroup v2
	}
	return 1 // cgroup v1
}

func round(v float64, digits int) float64 {
	pow := math.Pow(10, float64(digits))
	return math.Round(v*pow) / pow
}

func formatBytes(b uint64) (float64, string) {
	switch {
	case b >= GiB:
		return float64(b) / float64(GiB), "GiB"
	case b >= MiB:
		return float64(b) / float64(MiB), "MiB"
	case b >= KiB:
		return float64(b) / float64(KiB), "KiB"
	default:
		return float64(b), "B"
	}
}
