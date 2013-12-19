package proto

import (
	"io"
)

type Envelop struct {
	Header  *Header
	Payload Message
}

type Message interface {
	GetTypeId() byte
	Write(w io.Writer) error
}
