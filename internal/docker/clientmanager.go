package docker

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/moby/moby/client"
)

type HostConfig struct {
	Name string
	Addr string
}

// const certPath = "../certs/"
var certPath = ""

// 원격지 컨테이너 클라이언트 관리
type DockerClientManager struct {
	mu      sync.RWMutex
	clients map[string]*Client // key : container server host name
}

func SetCertpaht(path string) {
	certPath = path
}

/*
	manager, _ := NewDockerClientManager([]HostConfig{
		{Name: "local", Addr: ""},
		{Name: "prod-1", Addr: "tcp://10.0.0.10:2376"},
		{Name: "prod-2", Addr: "tcp://10.0.0.11:2376"},
	})
*/
func NewDockerClientManager(hosts []HostConfig) (*DockerClientManager, error) {
	m := &DockerClientManager{
		clients: make(map[string]*Client),
	}

	for _, h := range hosts {
		// raw, err := newSDKClient(h.Addr)
		raw, err := newSDKClientTLS(h.Addr)
		if err != nil {
			return nil, err
		}

		m.clients[h.Name] = &Client{
			cli:  raw,
			addr: h.Addr,
			name: h.Name,
		}
	}

	return m, nil
}

func newSDKClient(addr string) (*client.Client, error) {
	opts := []client.Opt{
		client.WithAPIVersionNegotiation(),
	}

	switch addr {
	case "", "unix":
		opts = append(opts,
			client.WithHost("unix:///var/run/docker.sock"),
		)
	default:
		opts = append(opts,
			client.WithHost(addr),
		)
	}

	return client.NewClientWithOpts(opts...)
}

func newSDKClientTLS(addr string) (*client.Client, error) {
	opts := []client.Opt{
		client.WithAPIVersionNegotiation(),
		client.WithHost(addr),
		client.WithTLSClientConfig(
			filepath.Join(certPath, "ca.pem"),
			filepath.Join(certPath, "cert.pem"),
			filepath.Join(certPath, "key.pem"),
		),
	}

	return client.NewClientWithOpts(opts...)
}

// need tset!
func NewDockerClientTLS(addr, cert, key, ca string) (*client.Client, error) {
	return client.NewClientWithOpts(
		client.WithHost(addr),
		client.WithTLSClientConfig(ca, cert, key),
		client.WithAPIVersionNegotiation(),
	)
}

func (m *DockerClientManager) Get(name string) (*Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	c, ok := m.clients[name]
	if !ok {
		return nil, fmt.Errorf("docker host not found: %s", name)
	}
	return c, nil
}

func (m *DockerClientManager) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for addr, c := range m.clients {
		_ = c.cli.Close()
		delete(m.clients, addr)
	}
}
