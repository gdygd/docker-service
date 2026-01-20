package service

import (
	"context"

	"docker_service/internal/docker"
)

type ServiceInterface interface {
	Test()
	ContainerList(ctx context.Context) ([]docker.Container, error)
	InspectContainer(ctx context.Context, containerID string) (docker.ContainerInspect, error)
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) error
	ContainerStats(ctx context.Context, id string, stream bool) (*docker.ContainerStats, error)
}
