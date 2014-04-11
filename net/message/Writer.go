package message

import (
	"encoding/binary"
	"io"
	"bytes"
)

type Writer interface {
	Write(m *Message) (err error)
}

type writer struct{
	w io.Writer
}

func NewWriter(w io.Writer) Writer {
	return &writer{
		w: w,
	}
}

func (b *writer) Write(m *Message) (err error) {
	defer func() {
		if err != nil {
			log.Debug("error writing message %v: %v", m, err)
		} else {
			log.Debug("message written: %v", m)
		}
	}()

	var buf bytes.Buffer

	length := uint16(len(m.Body))
	if err = binary.Write(&buf, ByteOrder, MagicStart); err != nil {
		log.Debug("error writing magic start %v: %v", MagicStart, err)
		return
	}
	if err = binary.Write(&buf, ByteOrder, m.TypeId); err != nil {
		log.Debug("error writing type id %v: %v", m.TypeId, err)
		return
	}
	if err = binary.Write(&buf, ByteOrder, m.CorrelationId); err != nil {
		log.Debug("error writing correlation id %v: %v", m.CorrelationId, err)
		return
	}
	if err = binary.Write(&buf, ByteOrder, length); err != nil {
		log.Debug("error writing length %v: %v", length, err)
		return
	}
	if err = binary.Write(&buf, ByteOrder, m.Body); err != nil {
		log.Debug("error writing body: %v", string(m.Body))
		return
	}

	if _, err = buf.WriteTo(b.w); err != nil {
		log.Debug("error flushing writer: %v", err)
		return
	}

	return
}
