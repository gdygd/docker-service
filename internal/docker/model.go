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

type ContainerInspect struct {
	ID    string
	Image string
	Name  string
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
