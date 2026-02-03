package collector

import (
	"docker_service/internal/logger"
	"docker_service/internal/pipeline"
	"sync"
)

// RingBuffer 버퍼가 가득 차면 오래된 데이터를 제거하는 채널 래퍼
type RingBuffer struct {
	ch   chan pipeline.Message
	size int
	mu   sync.Mutex
}

// NewRingBuffer RingBuffer 생성
func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		ch:   make(chan pipeline.Message, size),
		size: size,
	}
}

// Send 메시지 전송 (버퍼 full이면 오래된 데이터 제거 후 추가)
func (rb *RingBuffer) Send(msg pipeline.Message) {
	rb.mu.Lock()
	defer rb.mu.Unlock()

	select {
	case rb.ch <- msg:
		// 정상 전송
	default:
		// 버퍼 full - 오래된 데이터 제거 후 추가
		select {
		case dropped := <-rb.ch:
			logger.Log.Print(1, "[RingBuffer] buffer full, dropped old message: type=%s host=%s",
				dropped.Type, dropped.Host)
		default:
			// 채널이 비어있는 경우 (race condition 방지)
		}

		// 새 메시지 추가
		select {
		case rb.ch <- msg:
		default:
			logger.Log.Error("[RingBuffer] failed to send message after drop")
		}
	}
}

// Channel 읽기용 채널 반환
func (rb *RingBuffer) Channel() <-chan pipeline.Message {
	return rb.ch
}

// Close 채널 닫기
func (rb *RingBuffer) Close() {
	close(rb.ch)
}

// Len 현재 버퍼에 있는 메시지 수
func (rb *RingBuffer) Len() int {
	return len(rb.ch)
}
