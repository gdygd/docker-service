package docker

import (
	"context"

	"github.com/moby/moby/client"
)

// container 관련 API

// container list
func (c *Client) ListContainers(ctx context.Context) ([]Container, error) {
	results, err := c.cli.ContainerList(ctx, client.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	result := make([]Container, 0, len(results.Items))
	for _, v := range results.Items {
		result = append(result, Container{
			ID:     v.ID[:12],
			Name:   v.Names[0][1:],
			Image:  v.Image,
			State:  string(v.State),
			Status: v.Status,
		})
	}
	return result, nil
}

// 컨테이너 상세조회
//InspectContainer()
/*
Env
Port Binding
Volume Mount
Network
Restart Policy
Resource Limit
*/
func (c *Client) InspectContainer(ctx context.Context, containerID string) (client.ContainerInspectResult, error) {
	return c.cli.ContainerInspect(ctx, containerID, client.ContainerInspectOptions{})
}

/*
start
ContainerStart()
stop
ContainerStop()
restart
ContainerRestart()
pause & unpause
ContainerPause()
remove
ContainerRemove()
*/

func (c *Client) StartContainer(ctx context.Context, id string) (client.ContainerStartResult, error) {
	return c.cli.ContainerStart(ctx, id, client.ContainerStartOptions{})
}

func (c *Client) StopContainer(ctx context.Context, id string) (client.ContainerStopResult, error) {
	return c.cli.ContainerStop(ctx, id, client.ContainerStopOptions{})
}

// container resource monitoring
/*
CPU 사용률
Memory 사용량 / Limit
Network Rx / Tx
Block IO
*/
// ContainerStats()
func (c *Client) ContainerStats(ctx context.Context, id string, stream bool) (client.ContainerStatsResult, error) {
	return c.cli.ContainerStats(ctx, id, client.ContainerStatsOptions{Stream: stream})
}

// log
// ContainerLogs()

// event stream
/*
컨테이너 생성/종료 감지

장애 감지

Auto-restart 트리거
*/
