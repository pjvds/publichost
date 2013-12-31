package stream

import (
	"sync"
)

type threadSafeMap struct {
	lock    sync.Locker
	streams map[Id]Stream
}

func NewThreadSafeMap() Map {
	return &threadSafeMap{
		lock:    &sync.Mutex{},
		streams: make(map[Id]Stream),
	}
}

func (m *threadSafeMap) Add(id Id, s Stream) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.streams[id]; ok {
		return ErrAlreadyExists
	}

	m.streams[id] = s
	return nil
}

func (m *threadSafeMap) Get(id Id) (s Stream, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if v, ok := m.streams[id]; ok {
		s = v
	} else {
		err = ErrNotFound
	}

	return
}

func (m *threadSafeMap) Delete(id Id) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.streams[id]; !ok {
		return ErrNotFound
	}

	delete(m.streams, id)

	return nil
}
