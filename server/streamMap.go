package server

import (
	"errors"
	"sync"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
)

// A thread-safe map that holds streams.
type streamMap struct {
	streams map[uint32]*Stream
	lock    sync.Mutex
}

func newStreamMap() *streamMap {
	return &streamMap{
		streams: make(map[uint32]*Stream),
	}
}

func (s *streamMap) Add(streamId uint32, stream *Stream) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.streams[streamId]; ok {
		return ErrAlreadyExists
	}

	s.streams[streamId] = stream
	return nil
}

func (s *streamMap) Get(streamId uint32) (stream *Stream, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if v, ok := s.streams[streamId]; ok {
		stream = v
	} else {
		err = ErrNotFound
	}

	return
}

func (s *streamMap) Delete(streamId uint32) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.streams[streamId]; ok {
		delete(s.streams, streamId)
	} else {
		err = ErrNotFound
	}

	return
}
