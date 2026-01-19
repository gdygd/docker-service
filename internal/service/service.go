package service

import (
	"context"

	"docker_service/internal/docker"
)

type ServiceInterface interface {
	Test()
	ContainerList(ctx context.Context) ([]docker.Container, error)
	InspectContainer(ctx context.Context, containerID string) (docker.ContainerInspect, error)
}
