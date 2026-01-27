package docker

import (
	"context"

	"docker_service/internal/logger"

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

// func (c *Client) EventStream(ctx context.Context) client.EventsResult {
// 	stream := c.cli.Events(ctx, client.EventsListOptions{})

// 	for {
// 		select {
// 		case msg, ok := <-stream.Messages:
// 			if !ok {
// 				logger.Log.Error("event channel closed ")
// 			} else {
// 				// logger.Log.Print(2, "event :%v", msg)
// 				var evtmsg events.Message
// 				evtmsg = msg

// 				logger.Log.Print(2, "type : %s, action : %s, tm: %v", evtmsg.Type, evtmsg.Action, evtmsg.Time)
// 				logger.Log.Print(2, "Actor:%v", evtmsg.Actor.ID)

// 				keylist := []string{
// 					"image", "name", "exitCode", "execDuration",
// 				}
// 				for k, v := range evtmsg.Actor.Attributes {
// 					for _, kval := range keylist {
// 						if k == kval {
// 							logger.Log.Print(2, "(%v) [%v]", k, v)
// 						}
// 					}
// 				}

// 			}
// 		case <-ctx.Done():
// 			logger.Log.Print(2, "exit event stream..by ctx")
// 			return client.EventsResult{}
// 		}
// 	}
// }

// EventStreamRaw는 Docker Events API의 결과를 직접 반환 (EventManager용)
func (c *Client) EventStreamRaw(ctx context.Context) client.EventsResult {
	return c.cli.Events(ctx, client.EventsListOptions{})
}

func (c *Client) EventStream(ctx context.Context) client.EventsResult {
	// 이벤트 액션 맵 초기화
	InitEventAction()

	stream := c.cli.Events(ctx, client.EventsListOptions{})

	for {
		select {
		case msg, ok := <-stream.Messages:
			if !ok {
				logger.Log.Error("event channel closed")
				return client.EventsResult{}
			}

			evt := msg
			evtType := string(evt.Type)
			evtAction := string(evt.Action)

			// 필터 함수 사용
			if !FilterEvent(evtType, evtAction) {
				continue
			}

			// 기본 이벤트 로그
			logger.Log.Print(2, "[EVENT] type=%s action=%s time=%d id=%s", evtType, evtAction, evt.Time, evt.Actor.ID)

			// Attribute 화이트리스트 출력
			for _, key := range EvtAttribytes {
				if v, ok := evt.Actor.Attributes[key]; ok {
					logger.Log.Print(2, "  - %s = %s", key, v)
				}
			}

		case <-ctx.Done():
			logger.Log.Print(2, "exit event stream..by ctx")
			return client.EventsResult{}
		}
	}
}

// func (c *Client) WatchContainerEvents(
// 	ctx context.Context,
// 	eventCh chan<- events.Message,
// 	errCh chan<- error,
// ) {
// 	// container 이벤트만 필터링
// 	args := filters.NewArgs()
// 	args.Add("type", string(events.ContainerEventType))

// 	// events API 호출
// 	msgCh, errs := c.cli.Events(ctx, client.EventsListOptions{
// 		Filters: args,
// 	})

// 	for {
// 		select {
// 		case msg, ok := <-msgCh:
// 			if !ok {
// 				errCh <- fmt.Errorf("events channel closed")
// 				return
// 			}
// 			eventCh <- msg

// 		case err, ok := <-errs:
// 			if !ok {
// 				return
// 			}
// 			if err != nil {
// 				errCh <- err
// 				return
// 			}

// 		case <-ctx.Done():
// 			return
// 		}
// 	}
// }
