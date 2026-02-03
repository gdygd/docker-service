package collector

import (
	"context"

	"docker_service/internal/pipeline"
)

// Collector 데이터 수집 인터페이스
type Collector interface {
	// Name 수집기 이름
	Name() string

	// Start 수집 시작 (채널로 메시지 전송)
	Start(ctx context.Context) (<-chan pipeline.Message, error)

	// Stop 수집 중지
	Stop() error
}

// Config 수집기 공통 설정
type Config struct {
	// Host Docker 호스트명
	Host string

	// Interval 수집 주기 (List, Stats용)
	IntervalSec int

	// BufferSize 채널 버퍼 크기
	BufferSize int
}

// DefaultConfig 기본 설정
func DefaultConfig(host string) Config {
	return Config{
		Host:        host,
		IntervalSec: 10,
		BufferSize:  100,
	}
}
