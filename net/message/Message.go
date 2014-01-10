package message

import (
	"encoding/binary"
)

const (
	OpOpenTunnel  = byte(iota)
	OpOpenStream  = byte(iota)
	OpReadStream  = byte(iota)
	OpWriteStream = byte(iota)
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

func NewMessage(typeId byte, correlationId uint64, body []byte) *Message {
	return &Message{
		TypeId:        typeId,
		CorrelationId: correlationId,
		Body:          body,
	}
}
