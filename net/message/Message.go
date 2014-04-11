package message

import (
	"encoding/binary"
	"fmt"
)

const (
	OpPing        = byte(iota)
	OpOpenTunnel  = byte(iota)
	OpOpenStream  = byte(iota)
	OpReadStream  = byte(iota)
	OpWriteStream = byte(iota)
	OpCloseStream = byte(iota)
	Ack           = byte(iota)
	Nack          = byte(iota)
)

const (
	MagicStart = byte(42)
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

func (m *Message) String() string {
	return fmt.Sprintf("%v %v", GetTypeName(m.TypeId), m.CorrelationId)
}

func GetTypeName(t byte) string {
	switch t {
	case OpPing:
		return "ping"
	case OpOpenTunnel:
		return "open-tunnel"
	case OpOpenStream:
		return "open-stream"
	case OpReadStream:
		return "read-stream"
	case OpWriteStream:
		return "write-stream"
	case OpCloseStream:
		return "close-stream"
	case Ack:
		return "ack"
	case Nack:
		return "nack"
	default:
		return fmt.Sprintf("unknown(%v)", t)
	}
}