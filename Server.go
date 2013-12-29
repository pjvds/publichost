package publichost

import (
	"encoding/binary"
	"io"
)

type Message struct {
	Id            int16  // The id of the message
	TypeId        byte   // The type of the message
	CorrelationId int16  // The correlation id of the message
	Data          []byte // The data of the message
}

func (m *Message) WriteTo(writer io.Writer) (n int64, err error) {
	binary.Write
}
