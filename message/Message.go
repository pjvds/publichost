package message

import (
	"io"
)

type Message interface {
	GetTypeId() byte
	Write(w io.Writer) error
}
