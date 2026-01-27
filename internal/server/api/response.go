package api

import (
    "fmt"
    "math"

    "docker_service/internal/docker"
)

// ============================================================================
// Generic API Response
// ============================================================================

type APIResponse[T any] struct {
    Success bool   `json:"success"`
    Message string `json:"message,omitempty"`
    Data    T      `json:"data,omitempty"`
}

func SuccessResponse[T any](data T) APIResponse[T] {
    return APIResponse[T]{
        Success: true,
        Data:    data,
    }
}

func SuccessMessageResponse[T any](message string, data T) APIResponse[T] {
    return APIResponse[T]{
        Success: true,
        Message: message,
        Data:    data,
    }
}

func ErrorResponse(message string) APIResponse[any] {
    return APIResponse[any]{
        Success: false,
        Message: message,
    }
}

// ============================================================================
// Container List Response
// ============================================================================

type ContainerResponse struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    Image  string `json:"image"`
    State  string `json:"state"`
    Status string `json:"status"`
}

func ToContainerResponse(c docker.Container) ContainerResponse {
    return ContainerResponse{
        ID:     c.ID,
        Name:   c.Name,
        Image:  c.Image,
        State:  c.State,
        Status: c.Status,
    }
}

func ToContainerListResponse(containers []docker.Container) []ContainerResponse {
    result := make([]ContainerResponse, 0, len(containers))
    for _, c := range containers {
        result = append(result, ToContainerResponse(c))
    }
    return result
}

// ============================================================================
// Container Inspect Response
// ============================================================================

type ContainerInspectResponse struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Image string `json:"image"`
}

func ToContainerInspectResponse(c docker.ContainerInspect) ContainerInspectResponse {
    return ContainerInspectResponse{
        ID:    c.ID,
        Name:  c.Name,
        Image: c.Image,
    }
}

// ============================================================================
// Container Stats Response
// ============================================================================

type ContainerStatsResponse struct {
    ID            string  `json:"id"`
    Name          string  `json:"name"`
    CPUPercent    float64 `json:"cpu_percent"`
    MemoryUsage   string  `json:"memory_usage"`   // "1.2 GiB"
    MemoryLimit   string  `json:"memory_limit"`   // "4.0 GiB"
    MemoryPercent float64 `json:"memory_percent"`
    NetworkRx     string  `json:"network_rx"`     // "1.5 MiB"
    NetworkTx     string  `json:"network_tx"`     // "2.3 MiB"
}

func ToContainerStatsResponse(s *docker.ContainerStats) ContainerStatsResponse {
    return ContainerStatsResponse{
        ID:            s.ID,
        Name:          s.Name,
        CPUPercent:    roundFloat(s.CPUPercent, 2),
        MemoryUsage:   fmt.Sprintf("%.2f %s", s.MemoryUsageVal, s.MemoryUsageUnit),
        MemoryLimit:   fmt.Sprintf("%.2f %s", s.MemoryLimitVal, s.MemoryLimitUnit),
        MemoryPercent: roundFloat(s.MemoryPercent, 2),
        NetworkRx:     formatBytes(s.NetworkRx),
        NetworkTx:     formatBytes(s.NetworkTx),
    }
}

// ============================================================================
// Helper Functions
// ============================================================================

const (
    _          = iota
    KB float64 = 1 << (10 * iota)
    MB
    GB
    TB
)

func formatBytes(bytes uint64) string {
    b := float64(bytes)
    switch {
    case b >= TB:
        return fmt.Sprintf("%.2f TB", b/TB)
    case b >= GB:
        return fmt.Sprintf("%.2f GB", b/GB)
    case b >= MB:
        return fmt.Sprintf("%.2f MB", b/MB)
    case b >= KB:
        return fmt.Sprintf("%.2f KB", b/KB)
    default:
        return fmt.Sprintf("%d B", bytes)
    }
}

func roundFloat(val float64, precision int) float64 {
    ratio := math.Pow(10, float64(precision))
    return math.Round(val*ratio) / ratio
}
