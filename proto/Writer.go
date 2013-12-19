package proto

import (
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

func (r *Writer) Write(e Envelop) (err error) {
	encoder := gob.NewEncoder(r.writer)

	if err := encoder.Encode(e.Header); err != nil {
		return err
	}

	if err := e.Payload.Write(r.writer); err != nil {
		return err
	}

	return nil
}
