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
	StartContainer(ctx context.Context, id string) (client.ContainerStartResult, error)
	StopContainer(ctx context.Context, id string) (client.ContainerStopResult, error)
	ContainerPause(ctx context.Context, id string) error
	ContainerRemove(ctx context.Context, id string) error
	ContainerStats(ctx context.Context, id string, stream bool) (client.ContainerStatsResult, error)

	EventStream(ctx context.Context) client.EventsResult
	EventStreamRaw(ctx context.Context) client.EventsResult
}

// Docker Host
type Client struct {
	cli  *client.Client
	addr string
	name string // host name
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

func (c *Client) Addr() string {
	return c.addr
}

func (c *Client) Name() string {
	return c.name
}

func (c *Client) Raw() *client.Client {
	return c.cli
}
