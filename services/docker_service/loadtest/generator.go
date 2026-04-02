package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"docker_service/internal/pipeline"
)

var eventActions = [...]string{"start", "stop", "die", "restart", "kill", "pause", "unpause"}

// Generator produces synthetic pipeline.Message values for a single agent.
// It cycles through all four DataTypes in sequence.
type Generator struct {
	agentId    int
	containers int
	seq        int
	rng        *rand.Rand
}

func NewGenerator(agentId, containers int) *Generator {
	return &Generator{
		agentId:    agentId,
		containers: containers,
		rng:        rand.New(rand.NewPCG(uint64(agentId), 0)),
	}
}

// Next returns the next synthetic message, cycling list → inspect → stats → event.
func (g *Generator) Next() pipeline.Message {
	idx := g.seq % 4
	g.seq++

	msg := pipeline.Message{
		AgentId:   g.agentId,
		Host:      fmt.Sprintf("host-%04d", g.agentId),
		Timestamp: time.Now(),
	}

	switch idx {
	case 0:
		msg.Type = pipeline.DataTypeList
		msg.Data = g.genList()
	case 1:
		msg.Type = pipeline.DataTypeInspect
		msg.Data = g.genInspect()
	case 2:
		msg.Type = pipeline.DataTypeStats
		msg.Data = g.genStats()
	case 3:
		msg.Type = pipeline.DataTypeEvent
		msg.Data = g.genEvent()
	}
	return msg
}

func (g *Generator) genList() pipeline.ContainerListData {
	containers := make([]pipeline.ContainerInfo, g.containers)
	for i := range containers {
		containers[i] = pipeline.ContainerInfo{
			ID:     g.containerID(i),
			Name:   fmt.Sprintf("ct-%04d-%02d", g.agentId, i),
			Image:  "nginx:1.27",
			State:  "running",
			Status: "Up 2 hours",
		}
	}
	return pipeline.ContainerListData{Containers: containers}
}

func (g *Generator) genInspect() pipeline.ContainerInspectData {
	inspects := make([]pipeline.ContainerInspectInfo, g.containers)
	for i := range inspects {
		id := g.containerID(i)
		now := time.Now()
		inspects[i] = pipeline.ContainerInspectInfo{
			ID:       id[:12],
			ID2:      id,
			Name:     fmt.Sprintf("ct-%04d-%02d", g.agentId, i),
			Image:    "nginx:1.27",
			Created:  now.Add(-2 * time.Hour).Format(time.RFC3339),
			Platform: "linux",
			State: &pipeline.ContainerStateInfo{
				Status:    "running",
				Running:   true,
				ExitCode:  0,
				StartedAt: now.Add(-2 * time.Hour).Format(time.RFC3339),
			},
			Config: &pipeline.ContainerConfigInfo{
				Hostname: fmt.Sprintf("ct-%04d-%02d", g.agentId, i),
				Env:      []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
				Cmd:      []string{"nginx", "-g", "daemon off;"},
			},
			Network: &pipeline.ContainerNetworkInfo{
				IPAddress:  fmt.Sprintf("172.17.%d.%d", g.agentId/256, i+2),
				Gateway:    "172.17.0.1",
				MacAddress: fmt.Sprintf("02:42:ac:11:%02x:%02x", g.agentId%256, i+2),
				Ports:      map[string][]pipeline.PortBindingInfo{},
				Networks:   map[string]pipeline.NetworkEndpoint{},
			},
		}
	}
	return pipeline.ContainerInspectData{Inspects: inspects}
}

func (g *Generator) genStats() pipeline.ContainerStatsData {
	stats := make([]pipeline.ContainerStatsInfo, g.containers)
	for i := range stats {
		memUsage := uint64(g.rng.IntN(512)) * 1024 * 1024
		stats[i] = pipeline.ContainerStatsInfo{
			ID:            g.containerID(i),
			Name:          fmt.Sprintf("ct-%04d-%02d", g.agentId, i),
			CPUPercent:    g.rng.Float64() * 80,
			MemoryUsage:   memUsage,
			MemoryLimit:   2 * 1024 * 1024 * 1024,
			MemoryPercent: float64(memUsage) / float64(2*1024*1024*1024) * 100,
			NetworkRx:     uint64(g.rng.IntN(1024)) * 1024,
			NetworkTx:     uint64(g.rng.IntN(1024)) * 1024,
		}
	}
	return pipeline.ContainerStatsData{Stats: stats}
}

func (g *Generator) genEvent() pipeline.ContainerEvent {
	action := eventActions[g.rng.IntN(len(eventActions))]
	ctIdx := g.rng.IntN(g.containers)
	return pipeline.ContainerEvent{
		Host:      fmt.Sprintf("host-%04d", g.agentId),
		Type:      "container",
		Action:    action,
		ActorID:   g.containerID(ctIdx)[:12],
		ActorName: fmt.Sprintf("ct-%04d-%02d", g.agentId, ctIdx),
		Timestamp: time.Now().UnixMilli(),
		Attrs:     map[string]string{"image": "nginx:1.27"},
	}
}

// containerID generates a deterministic-looking 64-char hex container ID.
func (g *Generator) containerID(ctIdx int) string {
	// return fmt.Sprintf("%016x%016x%016x%016x",
	// 	uint64(g.agentId)<<32|uint64(ctIdx),
	// 	g.rng.Uint64(),
	// 	g.rng.Uint64(),
	// 	g.rng.Uint64(),
	// )
	return fmt.Sprintf("%016x",
		uint64(g.agentId)<<32|uint64(ctIdx),
	)
}
