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
