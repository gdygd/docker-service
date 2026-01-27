package event
/*
container event manager
-event stream
-event subscribe
-event dispatch

[Docker Daemon A] ──┐
                    │   WatchHost()
[Docker Daemon B] ──┼──────────────▶ [EventManager]
                    │                     │
[Docker Daemon C] ──┘                     │ dispatcher()
                                          │
                    ┌─────────────────────┼─────────────────────┐
                    │                     │                     │
                    ▼                     ▼                     ▼
            Subscribe("sse")      Subscribe("ws")      Subscribe("alert")
                    │                     │                     │
                    ▼                     ▼                     ▼
              [SSE Client]         [WS Client]          [Alert System]


              
[Docker Daemon]
       │
       ▼
EventStreamRaw() ──▶ client.EventsResult{Messages, Err}
       │
       ▼
streamEvents() ──▶ <-stream.Messages
       │
       ▼
em.eventChan ──▶ dispatcher() ──▶ broadcast()
       │
       ▼
sub.Events ──▶ bridgeEventsToSSE() ──▶ goglib.SendSSE()
       │
       ▼
/events (SSE) ──▶ 클라이언트

*/
import (
    "context"
    "sync"
    "time"
    "fmt"
    
    "docker_service/internal/docker"
    "docker_service/internal/logger"
)

// Subscriber는 이벤트를 받을 채널
type Subscriber struct {
    ID     string
    Events chan ContainerEvent
    Filter func(ContainerEvent) bool // optional filter
}

type EventManager struct {
    ctx        context.Context
    cancel     context.CancelFunc
    docMng     *docker.DockerClientManager
    
    // 이벤트 채널 (내부용 - 모든 호스트에서 fan-in)
    eventChan  chan ContainerEvent
    
    // 구독자 관리
    subscribers map[string]*Subscriber
    subMu       sync.RWMutex
    
    // 호스트별 watcher 관리
    watchers    map[string]context.CancelFunc
    watcherMu   sync.Mutex
    
    wg          sync.WaitGroup
}

func NewEventManager(ctx context.Context, docMng *docker.DockerClientManager) *EventManager {
    ctx, cancel := context.WithCancel(ctx)
    
    return &EventManager{
        ctx:         ctx,
        cancel:      cancel,
        docMng:      docMng,
        eventChan:   make(chan ContainerEvent, 100),
        subscribers: make(map[string]*Subscriber),
        watchers:    make(map[string]context.CancelFunc),
    }
}

// Start는 EventManager를 시작하고 dispatcher를 실행
func (em *EventManager) Start() {
    em.wg.Add(1)
    go em.dispatcher()
    
    logger.Log.Print(2, "[EventManager] Started")
}

// Stop은 모든 watcher와 dispatcher를 종료
func (em *EventManager) Stop() {
    em.cancel()

    // 모든 구독자 채널 닫기 (for range 루프 종료시킴)
    em.subMu.Lock()
    for id, sub := range em.subscribers {
        close(sub.Events)
        delete(em.subscribers, id)
        logger.Log.Print(2, "[EventManager] Closing subscriber: %s", id)
    }
    em.subMu.Unlock()

    // watcher들 종료 대기
    em.wg.Wait()

    // 내부 채널 닫기
    close(em.eventChan)

    logger.Log.Print(2, "[EventManager] Stopped")
}

// WatchHost는 특정 호스트의 이벤트 스트림을 시작
func (em *EventManager) WatchHost(host string) error {
    em.watcherMu.Lock()
    defer em.watcherMu.Unlock()
    
    // 이미 watch 중이면 skip
    if _, exists := em.watchers[host]; exists {
        return nil
    }
    
    client, err := em.docMng.Get(host)
    if err != nil {
        return err
    }
    
    watchCtx, watchCancel := context.WithCancel(em.ctx)
    em.watchers[host] = watchCancel
    
    em.wg.Add(1)
    go em.watchHostEvents(watchCtx, host, client)
    
    logger.Log.Print(2, "[EventManager] Started watching host: %s", host)
    return nil
}

// UnwatchHost는 특정 호스트의 이벤트 스트림을 중지
func (em *EventManager) UnwatchHost(host string) {
    em.watcherMu.Lock()
    defer em.watcherMu.Unlock()
    
    if cancel, exists := em.watchers[host]; exists {
        cancel()
        delete(em.watchers, host)
        logger.Log.Print(2, "[EventManager] Stopped watching host: %s", host)
    }
}

// Subscribe는 이벤트 구독자를 등록
func (em *EventManager) Subscribe(id string, bufferSize int, filter func(ContainerEvent) bool) *Subscriber {
    em.subMu.Lock()
    defer em.subMu.Unlock()
    
    sub := &Subscriber{
        ID:     id,
        Events: make(chan ContainerEvent, bufferSize),
        Filter: filter,
    }
    em.subscribers[id] = sub
    
    logger.Log.Print(2, "[EventManager] Subscriber added: %s", id)
    return sub
}

// Unsubscribe는 구독자를 제거
// Stop()에서 이미 닫혔을 수 있으므로 안전하게 처리
func (em *EventManager) Unsubscribe(id string) {
    em.subMu.Lock()
    defer em.subMu.Unlock()

    if sub, exists := em.subscribers[id]; exists {
        // Stop()에서 이미 닫혔을 수 있으므로 recover
        func() {
            defer func() { recover() }()
            close(sub.Events)
        }()
        delete(em.subscribers, id)
        logger.Log.Print(2, "[EventManager] Subscriber removed: %s", id)
    }
}

// dispatcher는 이벤트를 모든 구독자에게 분배
func (em *EventManager) dispatcher() {
    defer em.wg.Done()
    
    for {
        select {
        case <-em.ctx.Done():
            return
        case evt, ok := <-em.eventChan:
            if !ok {
                return
            }
            em.broadcast(evt)
        }
    }
}

func (em *EventManager) broadcast(evt ContainerEvent) {
    
    em.subMu.RLock()
    defer em.subMu.RUnlock()
    
    for _, sub := range em.subscribers {
        // 필터가 있으면 체크
        if sub.Filter != nil && !sub.Filter(evt) {
            continue
        }
        
        // non-blocking send
        select {
        case sub.Events <- evt:
            
        default:
            // 버퍼 가득 찼으면 skip (또는 로그)
            logger.Log.Warn("[EventManager] Subscriber %s buffer full, dropping event", sub.ID)
        }
    }
}

// watchHostEvents는 단일 호스트의 이벤트를 감시 (재연결 로직 포함)
func (em *EventManager) watchHostEvents(ctx context.Context, host string, client *docker.Client) {
    defer em.wg.Done()
    
    backoff := time.Second
    maxBackoff := 30 * time.Second
    
    for {
        select {
        case <-ctx.Done():
            return
        default:
        }
        
        err := em.streamEvents(ctx, host, client)
        
        if ctx.Err() != nil {
            return // context 취소됨
        }
        
        // 연결 끊김 - backoff 후 재시도
        logger.Log.Warn("[EventManager] Host %s stream disconnected: %v, retrying in %v", 
            host, err, backoff)
        
        select {
        case <-ctx.Done():
            return
        case <-time.After(backoff):
        }
        
        // exponential backoff
        backoff = backoff * 2
        if backoff > maxBackoff {
            backoff = maxBackoff
        }
    }
}

func (em *EventManager) streamEvents(ctx context.Context, host string, client *docker.Client) error {
    // EventStreamRaw는 블로킹하지 않고 EventsResult를 직접 반환
    stream := client.EventStreamRaw(ctx)

    logger.Log.Print(2, "[EventManager] streamEvents started for host: %s", host)

    for {
        select {
        case <-ctx.Done():
            logger.Log.Print(2, "[EventManager] streamEvents ctx done for host: %s", host)
            return ctx.Err()

        case err := <-stream.Err:
            logger.Log.Error("[EventManager] streamEvents error for host %s: %v", host, err)
            return err

        case msg, ok := <-stream.Messages:
            if !ok {
                return fmt.Errorf("event channel closed")
            }

            evtType := string(msg.Type)
            evtAction := string(msg.Action)

            // 이벤트 필터링 (허용되지 않은 Type/Action은 skip)
            if !docker.FilterEvent(evtType, evtAction) {
                continue
            }

            logger.Log.Print(2, "[EventManager] Received event: type=%s action=%s", evtType, evtAction)

            // Attribute 필터링
            filteredAttrs := docker.FilterAttrs(msg.Actor.Attributes)

            evt := ContainerEvent{
                Host:      host,
                Type:      evtType,
                Action:    evtAction,
                ActorID:   msg.Actor.ID,
                Timestamp: msg.Time,
                Attrs:     filteredAttrs,
            }

            if name, ok := msg.Actor.Attributes["name"]; ok {
                evt.ActorName = name
            }

            // fan-in: 이벤트를 중앙 채널로
            select {
            case em.eventChan <- evt:
                logger.Log.Print(2, "[EventManager] Event dispatched: %s/%s", evt.Type, evt.Action)
            case <-ctx.Done():
                return ctx.Err()
            }
        }
    }
}
