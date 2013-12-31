package message

import (
	"sync/atomic"
)

type IdSequence interface {
	Next() Id
}

type idSequence struct {
	lastId *uint32
}

func NewIdSequence() IdSequence {
	lastId := uint32(0)
	return &idSequence{
		lastId: &lastId,
	}
}

func (sequence *idSequence) Next() Id {
	id := atomic.AddUint32(sequence.lastId, 1)
	return Id(id)
}
