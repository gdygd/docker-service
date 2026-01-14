package service

import (
	"context"
	"fmt"

	"docker_service/internal/db"
	"docker_service/internal/docker"
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
