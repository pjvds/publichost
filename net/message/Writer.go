package message

import (
	"bufio"
	"encoding/binary"
	"io"
)

type Writer interface {
	Write(m *Message) (err error)
}

type bufferedWriter struct {
	writer *bufio.Writer
}

func NewWriter(r io.Writer) Writer {
	writer := bufio.NewWriter(r)
	return &bufferedWriter{
		writer: writer,
	}
}

func (b *bufferedWriter) Write(m *Message) (err error) {
	defer func() {
		if err != nil {
			log.Debug("error writing message %v: %v", m, err)
		} else {
			log.Debug("message written: %v", m)
		}
	}()

	length := uint16(len(m.Body))

	if err = binary.Write(b.writer, ByteOrder, m.TypeId); err != nil {
		return
	}
	if err = binary.Write(b.writer, ByteOrder, m.CorrelationId); err != nil {
		return
	}
	if err = binary.Write(b.writer, ByteOrder, length); err != nil {
		return
	}
	if err = binary.Write(b.writer, ByteOrder, m.Body); err != nil {
		return
	}

	// Perform the actual write by flushing.
	if err = b.writer.Flush(); err != nil {
		return
	}

	return
}
