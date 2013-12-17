package message

import (
	"bytes"
	"encoding/gob"
	"io"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

func (r *Writer) Write(m Message) (err error) {
	buffer := new(bytes.Buffer)
	if err := m.Write(buffer); err != nil {
		return err
	}

	header := NewHeader(m.GetTypeId(), int32(buffer.Len()))

	encoder := gob.NewEncoder(r.writer)
	if err := encoder.Encode(header); err != nil {
		return err
	}

	if _, err := buffer.WriteTo(r.writer); err != nil {
		return err
	}

	return nil
}
