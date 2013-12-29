package server

import (
	"errors"
	"io"
	"sync"
)

var (
	ErrStreamNotFound = errors.New("stream not found")
)

// Represents a single connection with a publichost client server.
// One server is created per accepted client connection.
type StreamManager struct {
	dialer StreamDialer

	streamIdSequence *idSequence
	streams          *streamMap
	streamsLock      sync.Mutex
}

func NewStreamManager(dialer StreamDialer) *StreamManager {
	return &StreamManager{
		dialer:           dialer,
		streamIdSequence: newIdSequence(),
		streams:          newStreamMap(),
	}
}

func (c *StreamManager) OpenStream(network, address string) (streamId uint32, err error) {
	var conn StreamConnection
	if conn, err = c.dialer.Dial(network, address); err != nil {
		return
	}

	c.streamsLock.Lock()
	defer c.streamsLock.Unlock()

	streamId = c.streamIdSequence.Next()
	stream := &Stream{
		Id:   streamId,
		conn: conn,
	}

	c.streams.Add(streamId, stream)
	return
}

func (c *StreamManager) ProxyData(streamId uint32, data []byte) (context *SendData, err error) {
	var stream *Stream
	if stream, err = c.streams.Get(streamId); err != nil {
		return
	}

	context, err = stream.Send(data)
	return
}

func (c *StreamManager) CloseStream(streamId uint32) (err error) {
	c.streamsLock.Lock()
	defer c.streamsLock.Unlock()

	var stream *Stream
	if stream, err = c.streams.Get(streamId); err != nil {
		return
	}

	stream.Close()
	return
}
