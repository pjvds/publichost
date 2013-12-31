package stream_test

import (
	"bytes"
	"github.com/pjvds/publichost/stream"
)

type NullStream struct {
	bytes.Buffer
}

func NewNullStream() stream.Stream {
	return &NullStream{
		Buffer: bytes.Buffer{},
	}
}

func (s *NullStream) Close() error {
	return nil
}
