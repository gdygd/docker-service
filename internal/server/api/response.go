package api

import (
	"fmt"
	"math"
	"time"

	"docker_service/internal/config"
	"docker_service/internal/db"
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

type userResponse struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt.Time,
		CreatedAt:         user.CreatedAt.Time,
	}
}

type loginUserResponse struct {
	SessionID             string       `json:"session_id"`
	AcessToken            string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	User                  userResponse `json:"user"`
}

type renewAccessTokenResponse struct {
	AcessToken           string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

type DockerHostConfig struct {
	Name string `json:"host"`
	Addr string `json:"addr"`
}

func ToContainerHostResponse(hostinfos []config.DockerHostConfig) []DockerHostConfig {
	var hosts []DockerHostConfig = []DockerHostConfig{}
	for _, host := range hostinfos {
		hosts = append(hosts, DockerHostConfig{Name: host.Name, Addr: host.Addr})
	}
	return hosts
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
	// 기본 정보
	ID           string `json:"id"`
	Name         string `json:"name"`
	Image        string `json:"image"`
	Created      string `json:"created"`
	Platform     string `json:"platform"`
	RestartCount int    `json:"restart_count"`

	// 상태 정보
	State *StateResponse `json:"state,omitempty"`

	// 설정 정보
	Config *ConfigResponse `json:"config,omitempty"`

	// 네트워크 정보
	Network *NetworkResponse `json:"network,omitempty"`

	// 마운트 정보
	Mounts []MountResponse `json:"mounts,omitempty"`
}

type StateResponse struct {
	Status     string `json:"status"`
	Running    bool   `json:"running"`
	Paused     bool   `json:"paused"`
	Restarting bool   `json:"restarting"`
	ExitCode   int    `json:"exit_code"`
	StartedAt  string `json:"started_at,omitempty"`
	FinishedAt string `json:"finished_at,omitempty"`
}

type ConfigResponse struct {
	Hostname   string            `json:"hostname,omitempty"`
	User       string            `json:"user,omitempty"`
	Env        []string          `json:"env,omitempty"`
	Cmd        []string          `json:"cmd,omitempty"`
	Entrypoint []string          `json:"entrypoint,omitempty"`
	WorkingDir string            `json:"working_dir,omitempty"`
	Labels     map[string]string `json:"labels,omitempty"`
}

type NetworkResponse struct {
	IPAddress  string                             `json:"ip_address"`
	Gateway    string                             `json:"gateway"`
	MacAddress string                             `json:"mac_address"`
	Ports      map[string][]PortResponse          `json:"ports,omitempty"`
	Networks   map[string]NetworkEndpointResponse `json:"networks,omitempty"`
}

type PortResponse struct {
	HostIP   string `json:"host_ip"`
	HostPort string `json:"host_port"`
}

type NetworkEndpointResponse struct {
	NetworkID  string `json:"network_id"`
	IPAddress  string `json:"ip_address"`
	Gateway    string `json:"gateway"`
	MacAddress string `json:"mac_address"`
}

type MountResponse struct {
	Type        string `json:"type"`
	Name        string `json:"name,omitempty"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Mode        string `json:"mode"`
	ReadWrite   bool   `json:"rw"`
}

func ToContainerInspectResponse(c docker.ContainerInspect) ContainerInspectResponse {
	resp := ContainerInspectResponse{
		ID:           c.ID,
		Name:         c.Name,
		Image:        c.Image,
		Created:      c.Created,
		Platform:     c.Platform,
		RestartCount: c.RestartCount,
	}

	// State 변환
	if c.State != nil {
		resp.State = &StateResponse{
			Status:     c.State.Status,
			Running:    c.State.Running,
			Paused:     c.State.Paused,
			Restarting: c.State.Restarting,
			ExitCode:   c.State.ExitCode,
			StartedAt:  c.State.StartedAt,
			FinishedAt: c.State.FinishedAt,
		}
	}

	// Config 변환
	if c.Config != nil {
		resp.Config = &ConfigResponse{
			Hostname:   c.Config.Hostname,
			User:       c.Config.User,
			Env:        c.Config.Env,
			Cmd:        c.Config.Cmd,
			Entrypoint: c.Config.Entrypoint,
			WorkingDir: c.Config.WorkingDir,
			Labels:     c.Config.Labels,
		}
	}

	// Network 변환
	if c.NetworkSettings != nil {
		resp.Network = &NetworkResponse{
			IPAddress:  c.NetworkSettings.IPAddress,
			Gateway:    c.NetworkSettings.Gateway,
			MacAddress: c.NetworkSettings.MacAddress,
		}

		// Ports 변환
		if c.NetworkSettings.Ports != nil {
			resp.Network.Ports = make(map[string][]PortResponse)
			for port, bindings := range c.NetworkSettings.Ports {
				portBindings := make([]PortResponse, 0, len(bindings))
				for _, b := range bindings {
					portBindings = append(portBindings, PortResponse{
						HostIP:   b.HostIP,
						HostPort: b.HostPort,
					})
				}
				resp.Network.Ports[port] = portBindings
			}
		}

		// Networks 변환
		if c.NetworkSettings.Networks != nil {
			resp.Network.Networks = make(map[string]NetworkEndpointResponse)
			for name, ep := range c.NetworkSettings.Networks {
				resp.Network.Networks[name] = NetworkEndpointResponse{
					NetworkID:  ep.NetworkID,
					IPAddress:  ep.IPAddress,
					Gateway:    ep.Gateway,
					MacAddress: ep.MacAddress,
				}
			}
		}
	}

	// Mounts 변환
	if c.Mounts != nil {
		resp.Mounts = make([]MountResponse, 0, len(c.Mounts))
		for _, m := range c.Mounts {
			resp.Mounts = append(resp.Mounts, MountResponse{
				Type:        m.Type,
				Name:        m.Name,
				Source:      m.Source,
				Destination: m.Destination,
				Mode:        m.Mode,
				ReadWrite:   m.RW,
			})
		}
	}

	return resp
}

// ============================================================================
// Container Stats Response
// ============================================================================

type ContainerStatsResponse struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryUsage   string  `json:"memory_usage"` // "1.2 GiB"
	MemoryLimit   string  `json:"memory_limit"` // "4.0 GiB"
	MemoryPercent float64 `json:"memory_percent"`
	NetworkRx     string  `json:"network_rx"` // "1.5 MiB"
	NetworkTx     string  `json:"network_tx"` // "2.3 MiB"
}

func ToContainerStatsResponse(s docker.ContainerStats) ContainerStatsResponse {
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
