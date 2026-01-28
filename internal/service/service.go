package service

import (
	"context"

	"docker_service/internal/db"
	"docker_service/internal/docker"
)

type ServiceInterface interface {
	Test()
	ContainerList(ctx context.Context) ([]docker.Container, error)
	ContainerList2(ctx context.Context, host string) ([]docker.Container, error)
	InspectContainer(ctx context.Context, containerID string) (docker.ContainerInspect, error)
	InspectContainer2(ctx context.Context, containerID, host string) (docker.ContainerInspect, error)
	StartContainer(ctx context.Context, id string) error
	StartContainer2(ctx context.Context, id, host string) error
	StopContainer(ctx context.Context, id string) error
	StopContainer2(ctx context.Context, id, host string) error
	ContainerStats(ctx context.Context, id string, stream bool) (*docker.ContainerStats, error)
	ContainerStats2(ctx context.Context, id, host string, stream bool) (*docker.ContainerStats, error)
	ContainerStatsStream(ctx context.Context, id string, stream bool, ch_rst chan *docker.ContainerStats) error

	EventStream(ctx context.Context, host string)

	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	LoginUser(ctx context.Context, username string) (db.User, error)
	CreateSession(ctx context.Context, arg db.CreateSessionParams) (db.Session, error)
	ReadSession(ctx context.Context, id string) (db.Session, error)
}
