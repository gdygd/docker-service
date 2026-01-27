package docker

import (
    "github.com/moby/moby/client"
)

// ConvertInspectResult는 Docker SDK의 ContainerInspectResult를 내부 모델로 변환
func ConvertInspectResult(res client.ContainerInspectResult) ContainerInspect {
    c := res.Container

    inspect := ContainerInspect{
        ID:           c.ID,
        Name:         c.Name,
        Image:        c.Image,
        Created:      c.Created,
        Platform:     c.Platform,
        RestartCount: c.RestartCount,
    }

    // State 변환
    if c.State != nil {
        inspect.State = &ContainerState{
            Status:     string(c.State.Status),
            Running:    c.State.Running,
            Paused:     c.State.Paused,
            Restarting: c.State.Restarting,
            OOMKilled:  c.State.OOMKilled,
            Dead:       c.State.Dead,
            Pid:        c.State.Pid,
            ExitCode:   c.State.ExitCode,
            Error:      c.State.Error,
            StartedAt:  c.State.StartedAt,
            FinishedAt: c.State.FinishedAt,
        }
    }

    // Config 변환
    if c.Config != nil {
        inspect.Config = &ContainerConfig{
            Hostname:   c.Config.Hostname,
            User:       c.Config.User,
            Env:        c.Config.Env,
            Cmd:        c.Config.Cmd,
            WorkingDir: c.Config.WorkingDir,
            Labels:     c.Config.Labels,
        }

        // Entrypoint 변환
        if c.Config.Entrypoint != nil {
            inspect.Config.Entrypoint = c.Config.Entrypoint
        }

        // ExposedPorts 변환
        if c.Config.ExposedPorts != nil {
            inspect.Config.ExposedPorts = make(map[string]struct{})
            for port := range c.Config.ExposedPorts {
                inspect.Config.ExposedPorts[port.String()] = struct{}{}
            }
        }
    }

    // NetworkSettings 변환
    if c.NetworkSettings != nil {
        inspect.NetworkSettings = &ContainerNetworkSettings{}

        // Ports 변환
        if c.NetworkSettings.Ports != nil {
            inspect.NetworkSettings.Ports = make(map[string][]PortBinding)
            for port, bindings := range c.NetworkSettings.Ports {
                portKey := port.String()
                portBindings := make([]PortBinding, 0, len(bindings))
                for _, b := range bindings {
                    portBindings = append(portBindings, PortBinding{
                        HostIP:   b.HostIP.String(),
                        HostPort: b.HostPort,
                    })
                }
                inspect.NetworkSettings.Ports[portKey] = portBindings
            }
        }

        // Networks 변환 (IP/Gateway/MacAddress는 Networks 맵 안에 있음)
        if c.NetworkSettings.Networks != nil {
            inspect.NetworkSettings.Networks = make(map[string]NetworkEndpoint)
            for name, ep := range c.NetworkSettings.Networks {
                inspect.NetworkSettings.Networks[name] = NetworkEndpoint{
                    NetworkID:  ep.NetworkID,
                    IPAddress:  ep.IPAddress.String(),
                    Gateway:    ep.Gateway.String(),
                    MacAddress: ep.MacAddress.String(),
                }

                // 첫 번째 네트워크의 정보를 기본 IP/Gateway/MacAddress로 설정
                if inspect.NetworkSettings.IPAddress == "" {
                    inspect.NetworkSettings.IPAddress = ep.IPAddress.String()
                    inspect.NetworkSettings.Gateway = ep.Gateway.String()
                    inspect.NetworkSettings.MacAddress = ep.MacAddress.String()
                }
            }
        }
    }

    // Mounts 변환
    if c.Mounts != nil {
        inspect.Mounts = make([]MountPoint, 0, len(c.Mounts))
        for _, m := range c.Mounts {
            inspect.Mounts = append(inspect.Mounts, MountPoint{
                Type:        string(m.Type),
                Name:        m.Name,
                Source:      m.Source,
                Destination: m.Destination,
                Mode:        m.Mode,
                RW:          m.RW,
            })
        }
    }

    return inspect
}
