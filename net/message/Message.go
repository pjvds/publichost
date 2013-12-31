package message

import (
	"encoding/binary"
)

const (
	OpOpenStream  = byte(iota)
	OpStreamData  = byte(iota)
	OpCloseStream = byte(iota)
	Ack           = byte(iota)
	Nack          = byte(iota)
)

var (
	ByteOrder = binary.BigEndian
)

type Message struct {
	TypeId        byte
	CorrelationId uint64
	Body          []byte
}
