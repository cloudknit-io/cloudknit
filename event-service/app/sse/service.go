package sse

import (
	"github.com/pkg/errors"
)

type API interface {
	Close()
	Chan() chan any
	Send(message any) error
}

type SSE struct {
	pipe chan any
}

func NewSSE() *SSE {
	return &SSE{
		pipe: make(chan any),
	}
}

func (s *SSE) Close() {
	if s.pipe != nil {
		close(s.pipe)
		s.pipe = nil
	}
}

func (s *SSE) Chan() chan any {
	return s.pipe
}

func (s *SSE) Send(message any) error {
	if s.pipe == nil {
		return errors.New("cannot send message to a closed channel")
	}

	select {
	case s.pipe <- message:
		break
	default:
		break
	}

	return nil
}
