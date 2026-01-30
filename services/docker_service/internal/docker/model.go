package docker

// DTO
type Container struct {
	ID     string
	Name   string
	Image  string
	State  string
	Status string
}

type ContainerAction string

const (
	Start   ContainerAction = "start"
	Stop    ContainerAction = "stop"
	Restart ContainerAction = "restart"
)

// ============================================================================
// Container Inspect (기본 + 네트워크 + 설정)
// ============================================================================

type ContainerInspect struct {
    // 기본 정보
    ID           string
    Name         string
    Image        string
    Created      string
    Platform     string
    RestartCount int
    State        *ContainerState

    // 설정 정보
    Config *ContainerConfig

    // 네트워크 정보
    NetworkSettings *ContainerNetworkSettings

    // 마운트 정보
    Mounts []MountPoint
}

// ContainerState - 컨테이너 상태 정보
type ContainerState struct {
    Status     string // running, exited, paused, etc.
    Running    bool
    Paused     bool
    Restarting bool
    OOMKilled  bool
    Dead       bool
    Pid        int
    ExitCode   int
    Error      string
    StartedAt  string
    FinishedAt string
}

// ContainerConfig - 컨테이너 설정 정보
type ContainerConfig struct {
    Hostname     string
    User         string
    Env          []string
    Cmd          []string
    Entrypoint   []string
    WorkingDir   string
    ExposedPorts map[string]struct{}
    Labels       map[string]string
}

// ContainerNetworkSettings - 네트워크 설정
type ContainerNetworkSettings struct {
    IPAddress   string
    Gateway     string
    MacAddress  string
    Ports       map[string][]PortBinding
    Networks    map[string]NetworkEndpoint
}

// PortBinding - 포트 바인딩 정보
type PortBinding struct {
    HostIP   string
    HostPort string
}

// NetworkEndpoint - 네트워크 엔드포인트 정보
type NetworkEndpoint struct {
    NetworkID   string
    IPAddress   string
    Gateway     string
    MacAddress  string
}

// MountPoint - 마운트 정보
type MountPoint struct {
    Type        string // bind, volume, tmpfs
    Name        string // volume name (if volume)
    Source      string // host path
    Destination string // container path
    Mode        string // rw, ro
    RW          bool
}

type ContainerStatsRaw struct {
	CPUStats struct {
		CPUUsage struct {
			TotalUsage  uint64   `json:"total_usage"`
			PercpuUsage []uint64 `json:"percpu_usage"`
		} `json:"cpu_usage"`
		SystemCPUUsage uint64 `json:"system_cpu_usage"`
	} `json:"cpu_stats"`

	PreCPUStats struct {
		CPUUsage struct {
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemCPUUsage uint64 `json:"system_cpu_usage"`
	} `json:"precpu_stats"`

	// MemoryStats struct {
	// 	Usage uint64 `json:"usage"`
	// 	Limit uint64 `json:"limit"`
	// } `json:"memory_stats"`

	MemoryStats struct {
		Usage uint64 `json:"usage"`
		Limit uint64 `json:"limit"`
		Stats struct {
			Cache        uint64 `json:"cache"`         // cgroup v1
			InactiveFile uint64 `json:"inactive_file"` // cgroup v2
		} `json:"stats"`
	} `json:"memory_stats"`

	Networks map[string]struct {
		RxBytes uint64 `json:"rx_bytes"`
		TxBytes uint64 `json:"tx_bytes"`
	} `json:"networks"`
}

type ContainerStats struct {
	ID   string
	Name string

	CPUPercent  float64
	MemoryUsage uint64 // byte
	MemoryLimit uint64 // byte

	MemoryUsageVal  float64 // KiB ~ GiB 포캣적용 값
	MemoryUsageUnit string  // 포맷 단위 (KiB.. GiB)
	MemoryLimitVal  float64 // KiB ~ GiB 포캣적용 값
	MemoryLimitUnit string  // 포맷 단위 (KiB.. GiB)

	MemoryPercent float64
	NetworkRx     uint64 // byte
	NetworkTx     uint64 // byte
}
