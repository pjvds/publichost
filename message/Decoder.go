package message

import (
	"encoding/gob"
	"fmt"
	"io"
)

type Decoder interface {
	Decode(typeId byte, length int32, body io.Reader) (Message, error)
}

type hardcodedDecoder struct{}

func NewDecoder() Decoder {
	return &hardcodedDecoder{}
}

func (h *hardcodedDecoder) Decode(typeId byte, length int32, body io.Reader) (Message, error) {
	decoder := gob.NewDecoder(body)

	switch typeId {
	case TypeExposeResponse:
		msg := new(ExposeResponse)
		err := decoder.Decode(msg)
		return msg, err

	case TypeExposeRequest:
		msg := new(ExposeRequest)
		err := decoder.Decode(&msg)
		return msg, err
	}

	return nil, fmt.Errorf("Unknown type id %v", typeId)
}
