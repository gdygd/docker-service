package service

import (
	"context"
	"fmt"

	"docker_service/internal/db"
	"docker_service/internal/docker"
	"docker_service/internal/logger"
	"docker_service/internal/service"
)

type ApiService struct {
	dbHnd  db.DbHandler
	docker *docker.Client
}

func NewApiService(dbHnd db.DbHandler, docker *docker.Client) service.ServiceInterface {
	return &ApiService{
		dbHnd:  dbHnd,
		docker: docker,
	}
}

func (s *ApiService) Test() {
	fmt.Printf("test service")
}

func (s *ApiService) ContainerList(ctx context.Context) ([]docker.Container, error) {
	return s.docker.ListContainers(ctx)
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
