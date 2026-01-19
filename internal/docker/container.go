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

// container resource monitoring
/*
CPU 사용률
Memory 사용량 / Limit
Network Rx / Tx
Block IO
*/
// ContainerStats()

// log
// ContainerLogs()

// event stream
/*
컨테이너 생성/종료 감지

장애 감지

Auto-restart 트리거
*/
