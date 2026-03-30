package gapi

import (
	"fmt"

	"docker_service/internal/pipeline"
	"docker_service/pb"
)

// ConvertToAgentMessage converts pipeline.Message to pb.AgentMessage for DataStream transmission.
func ConvertToAgentMessage(msg pipeline.Message, agentKey string) (*pb.AgentMessage, error) {
	pbMsg := &pb.AgentMessage{
		Agentid:   int32(msg.AgentId),
		AgentKey:  agentKey,
		Type:      convertDataType(msg.Type),
		Host:      msg.Host,
		Timestamp: msg.Timestamp.UnixMilli(),
	}

	switch msg.Type {
	case pipeline.DataTypeList:
		data, err := assertData[pipeline.ContainerListData](msg.Data)
		if err != nil {
			return nil, fmt.Errorf("DataTypeList: %w", err)
		}
		pbMsg.Data = &pb.AgentMessage_ListData{ListData: convertListData(data)}

	case pipeline.DataTypeInspect:
		data, err := assertData[pipeline.ContainerInspectData](msg.Data)
		if err != nil {
			return nil, fmt.Errorf("DataTypeInspect: %w", err)
		}
		pbMsg.Data = &pb.AgentMessage_InspectData{InspectData: convertInspectData(data)}

	case pipeline.DataTypeStats:
		data, err := assertData[pipeline.ContainerStatsData](msg.Data)
		if err != nil {
			return nil, fmt.Errorf("DataTypeStats: %w", err)
		}
		pbMsg.Data = &pb.AgentMessage_StatsData{StatsData: convertStatsData(data)}

	case pipeline.DataTypeEvent:
		// pipeline.ContainerEventData is not yet defined; skip conversion
		// return nil, fmt.Errorf("DataTypeEvent: not yet implemented")
		data, err := assertData[pipeline.ContainerEvent](msg.Data)
		if err != nil {
			return nil, fmt.Errorf("DataTypeEvent: %w", err)
		}
		pbMsg.Data = &pb.AgentMessage_EventData{EventData: convertEventData(data)}

	default:
		return nil, fmt.Errorf("unknown DataType: %s", msg.Type)
	}

	return pbMsg, nil
}

func convertDataType(t pipeline.DataType) pb.DataType {
	switch t {
	case pipeline.DataTypeList:
		return pb.DataType_CONTAINER_LIST
	case pipeline.DataTypeInspect:
		return pb.DataType_CONTAINER_INSPECT
	case pipeline.DataTypeStats:
		return pb.DataType_CONTAINER_STATS
	case pipeline.DataTypeEvent:
		return pb.DataType_CONTAINER_EVENT
	default:
		return pb.DataType_CONTAINER_LIST
	}
}

// assertData performs type assertion supporting both value and pointer forms.
func assertData[T any](data interface{}) (T, error) {
	if v, ok := data.(T); ok {
		return v, nil
	}
	if p, ok := data.(*T); ok && p != nil {
		return *p, nil
	}
	var zero T
	return zero, fmt.Errorf("unexpected data type %T, want %T", data, zero)
}

// --- List ---

func convertListData(d pipeline.ContainerListData) *pb.ContainerListData {
	containers := make([]*pb.ContainerInfo, len(d.Containers))
	for i, c := range d.Containers {
		containers[i] = &pb.ContainerInfo{
			Id:     c.ID,
			Name:   c.Name,
			Image:  c.Image,
			State:  c.State,
			Status: c.Status,
		}
	}
	return &pb.ContainerListData{Containers: containers}
}

// --- Stats ---

func convertStatsData(d pipeline.ContainerStatsData) *pb.ContainerStatsData {
	stats := make([]*pb.ContainerStats, len(d.Stats))
	for i, s := range d.Stats {
		stats[i] = &pb.ContainerStats{
			Id:            s.ID,
			Name:          s.Name,
			CpuPercent:    s.CPUPercent,
			MemoryUsage:   s.MemoryUsage,
			MemoryLimit:   s.MemoryLimit,
			MemoryPercent: s.MemoryPercent,
			NetworkRx:     s.NetworkRx,
			NetworkTx:     s.NetworkTx,
		}
	}
	return &pb.ContainerStatsData{Stats: stats}
}

// --- Inspect ---

func convertInspectData(d pipeline.ContainerInspectData) *pb.ContainerInspectData {
	inspects := make([]*pb.ContainerInspect, len(d.Inspects))
	for i, ins := range d.Inspects {
		inspects[i] = &pb.ContainerInspect{
			Id:       ins.ID,
			Name:     ins.Name,
			Image:    ins.Image,
			Created:  ins.Created,
			Platform: ins.Platform,
			State:    convertContainerState(ins.State),
			Config:   convertContainerConfig(ins.Config),
			Network:  convertContainerNetwork(ins.Network),
			Mounts:   convertMounts(ins.Mounts),
		}
	}
	return &pb.ContainerInspectData{Inspects: inspects}
}

func convertContainerState(s *pipeline.ContainerStateInfo) *pb.ContainerState {
	if s == nil {
		return nil
	}
	return &pb.ContainerState{
		Status:     s.Status,
		Running:    s.Running,
		Paused:     s.Paused,
		Restarting: s.Restarting,
		ExitCode:   int32(s.ExitCode),
		StartedAt:  s.StartedAt,
		FinishedAt: s.FinishedAt,
	}
}

func convertContainerConfig(c *pipeline.ContainerConfigInfo) *pb.ContainerConfig {
	if c == nil {
		return nil
	}
	return &pb.ContainerConfig{
		Hostname:   c.Hostname,
		User:       c.User,
		Env:        c.Env,
		Cmd:        c.Cmd,
		Entrypoint: c.Entrypoint,
		WorkingDir: c.WorkingDir,
		Labels:     c.Labels,
	}
}

func convertContainerNetwork(n *pipeline.ContainerNetworkInfo) *pb.ContainerNetwork {
	if n == nil {
		return nil
	}

	ports := make(map[string]*pb.PortBindings, len(n.Ports))
	for k, bindings := range n.Ports {
		pbs := make([]*pb.PortBinding, len(bindings))
		for i, b := range bindings {
			pbs[i] = &pb.PortBinding{HostIp: b.HostIP, HostPort: b.HostPort}
		}
		ports[k] = &pb.PortBindings{Bindings: pbs}
	}

	networks := make(map[string]*pb.NetworkEndpoint, len(n.Networks))
	for k, ep := range n.Networks {
		networks[k] = &pb.NetworkEndpoint{
			NetworkId:  ep.NetworkID,
			IpAddress:  ep.IPAddress,
			Gateway:    ep.Gateway,
			MacAddress: ep.MacAddress,
		}
	}

	return &pb.ContainerNetwork{
		IpAddress:  n.IPAddress,
		Gateway:    n.Gateway,
		MacAddress: n.MacAddress,
		Ports:      ports,
		Networks:   networks,
	}
}

func convertMounts(mounts []pipeline.MountPointInfo) []*pb.MountPoint {
	result := make([]*pb.MountPoint, len(mounts))
	for i, m := range mounts {
		result[i] = &pb.MountPoint{
			Type:        m.Type,
			Name:        m.Name,
			Source:      m.Source,
			Destination: m.Destination,
			Mode:        m.Mode,
			Rw:          m.RW,
		}
	}
	return result
}

// --- event ---

func convertEventData(d pipeline.ContainerEvent) *pb.ContainerEventData {
	event := &pb.ContainerEventData{
		Type:      d.Type,
		Action:    d.Action,
		ActorId:   d.ActorID,
		ActorName: d.ActorName,
		Timestamp: d.Timestamp,
		Attrs:     d.Attrs,
	}

	return event
}
