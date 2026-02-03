package pipeline

import "time"

// DataType 전송 데이터 타입
type DataType string

const (
	DataTypeList    DataType = "container_list"
	DataTypeInspect DataType = "container_inspect"
	DataTypeStats   DataType = "container_stats"
	DataTypeEvent   DataType = "container_event"
)

// Message 파이프라인 통합 메시지 포맷
type Message struct {
	Type      DataType    `json:"type"`
	Host      string      `json:"host"` // Docker 호스트명
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// ContainerListData List 수집 데이터
type ContainerListData struct {
	Containers []ContainerInfo `json:"containers"`
}

type ContainerInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Image  string `json:"image"`
	State  string `json:"state"`
	Status string `json:"status"`
}
