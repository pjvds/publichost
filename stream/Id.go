package stream

import (
	"encoding/binary"
)

// The identifier of a stream.
// It ranges from 0 to 4,294,967,295.
type Id uint32

func (i Id) Bytes() []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))

	return b
}

func ParseId(body []byte) Id {
	id := binary.BigEndian.Uint32(body)
	return Id(id)
}