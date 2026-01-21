package api

import "time"

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type userResponse struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type requestContainerID struct {
	ID string `uri:"id" binding:"required"`
}

type ContainerStatsMessage struct {
	ContainerID string  `json:"containerId"`
	CPUPercent  float64 `json:"cpuPercent"`
	MemUsageMiB float64 `json:"memUsageMiB"`
	MemLimitMiB float64 `json:"memLimitMiB"`
	MemPercent  float64 `json:"memPercent"`
	NetRxBytes  uint64  `json:"netRxBytes"`
	NetTxBytes  uint64  `json:"netTxBytes"`
	BlockRead   uint64  `json:"blockRead"`
	BlockWrite  uint64  `json:"blockWrite"`
	Timestamp   int64   `json:"timestamp"`
}
