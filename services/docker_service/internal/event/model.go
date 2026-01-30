package event


// ContainerEvent는 정규화된 이벤트 구조체
type ContainerEvent struct {
    Host      string            `json:"host"`
    Type      string            `json:"type"`      // container, image, network...
    Action    string            `json:"action"`    // start, stop, die...
    ActorID   string            `json:"actor_id"`
    ActorName string            `json:"actor_name"`
    Timestamp int64             `json:"timestamp"`
    Attrs     map[string]string `json:"attrs,omitempty"`
}
