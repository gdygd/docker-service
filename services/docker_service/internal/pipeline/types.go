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

// ContainerInspectData Inspect 수집 데이터
type ContainerInspectData struct {
	Inspects []ContainerInspectInfo `json:"inspects"`
}

type ContainerInspectInfo struct {
	// 기본 정보
	ID           string `json:"id"`
	Name         string `json:"name"`
	Image        string `json:"image"`
	Created      string `json:"created"`
	Platform     string `json:"platform"`
	RestartCount int    `json:"restart_count"`

	// 상태 정보
	State *ContainerStateInfo `json:"state"`

	// 설정 정보
	Config *ContainerConfigInfo `json:"config"`

	// 네트워크 정보
	Network *ContainerNetworkInfo `json:"network"`

	// 마운트 정보
	Mounts []MountPointInfo `json:"mounts"`
}

type ContainerStateInfo struct {
	Status     string `json:"status"`
	Running    bool   `json:"running"`
	Paused     bool   `json:"paused"`
	Restarting bool   `json:"restarting"`
	ExitCode   int    `json:"exit_code"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
}

type ContainerConfigInfo struct {
	Hostname   string            `json:"hostname"`
	User       string            `json:"user"`
	Env        []string          `json:"env"`
	Cmd        []string          `json:"cmd"`
	Entrypoint []string          `json:"entrypoint"`
	WorkingDir string            `json:"working_dir"`
	Labels     map[string]string `json:"labels"`
}

type ContainerNetworkInfo struct {
	IPAddress  string                       `json:"ip_address"`
	Gateway    string                       `json:"gateway"`
	MacAddress string                       `json:"mac_address"`
	Ports      map[string][]PortBindingInfo `json:"ports"`
	Networks   map[string]NetworkEndpoint   `json:"networks"`
}

type PortBindingInfo struct {
	HostIP   string `json:"host_ip"`
	HostPort string `json:"host_port"`
}

type NetworkEndpoint struct {
	NetworkID  string `json:"network_id"`
	IPAddress  string `json:"ip_address"`
	Gateway    string `json:"gateway"`
	MacAddress string `json:"mac_address"`
}

type MountPointInfo struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
	RW          bool   `json:"rw"`
}

// ContainerStatsData Stats 수집 데이터
type ContainerStatsData struct {
	Stats []ContainerStatsInfo `json:"stats"`
}

type ContainerStatsInfo struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryUsage   uint64  `json:"memory_usage"`   // bytes
	MemoryLimit   uint64  `json:"memory_limit"`   // bytes
	MemoryPercent float64 `json:"memory_percent"`
	NetworkRx     uint64  `json:"network_rx"` // bytes
	NetworkTx     uint64  `json:"network_tx"` // bytes
}
