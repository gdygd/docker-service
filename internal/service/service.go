package service

import (
	"context"

	"docker_service/internal/docker"
)

type ServiceInterface interface {
	Test()
	ContainerList(ctx context.Context) ([]docker.Container, error)
	ContainerList2(ctx context.Context, host string) ([]docker.Container, error)
	InspectContainer(ctx context.Context, containerID string) (docker.ContainerInspect, error)
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) error
	ContainerStats(ctx context.Context, id string, stream bool) (*docker.ContainerStats, error)
	ContainerStatsStream(ctx context.Context, id string, stream bool, ch_rst chan *docker.ContainerStats) error
}
