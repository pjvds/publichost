package message

import (
	"encoding/gob"
	"fmt"
	"io"
)

type Decoder interface {
	Decode(typeId byte, length int32, body io.Reader) (interface{}, error)
}

type hardcodedDecoder struct{}

func NewDecoder() Decoder {
	return &hardcodedDecoder{}
}

func (h *hardcodedDecoder) Decode(typeId byte, length int32, body io.Reader) (interface{}, error) {
	decoder := gob.NewDecoder(body)

	switch typeId {
	case TypeExposeReponse:
		var msg ExposeResponse
		err := decoder.Decode(&msg)
		return msg, err

	case TypeExposeRequest:
		var msg ExposeRequest
		err := decoder.Decode(&msg)
		return msg, err
	}

	return nil, fmt.Errorf("Unknown type id %v", typeId)
}
