package server

import (
	"sync/atomic"
)

type idSequence struct {
	lastId *uint32
}

func newIdSequence() *idSequence {
	lastId := uint32(0)
	return &idSequence{
		lastId: &lastId,
	}
}

func (sequence *idSequence) Next() uint32 {
	return atomic.AddUint32(sequence.lastId, 1)
}
