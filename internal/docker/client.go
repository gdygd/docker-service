package docker

// Docker SDK 초기화
// ref : https://pkg.go.dev/github.com/moby/moby/client

import (
	"context"

	"github.com/moby/moby/client"
)

type DockerAPI interface {
	ListContainers(ctx context.Context) ([]Container, error)
	InspectContainer(ctx context.Context, containerID string) (client.ContainerInspectResult, error)
	StartContainer(ctx context.Context, id string) error
	StopContainer(ctx context.Context, id string) error
	ContainerPause(ctx context.Context, id string) error
	ContainerRemove(ctx context.Context, id string) error
	ContainerStats(ctx context.Context) error
}

type Client struct {
	cli *client.Client
}

func New() (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	return &Client{cli: cli}, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}
