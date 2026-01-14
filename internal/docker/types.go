package docker

// DTO
type Container struct {
	ID     string
	Name   string
	Image  string
	State  string
	Status string
}
