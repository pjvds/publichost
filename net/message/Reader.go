package message

import (
	"bufio"
	"encoding/binary"
	"io"
)

type Reader interface {
	Read() (m *Message, err error)
}

type bufferedReader struct {
	reader *bufio.Reader
}

func NewReader(r io.Reader) Reader {
	reader := bufio.NewReader(r)
	return &bufferedReader{
		reader: reader,
	}
}

func (b *bufferedReader) Read() (m *Message, err error) {
	var typeId byte
	var correlationId uint64
	var length uint16
	var body []byte

	if err = binary.Read(b.reader, ByteOrder, &typeId); err != nil {
		return
	}
	if err = binary.Read(b.reader, ByteOrder, &correlationId); err != nil {
		return
	}
	if err = binary.Read(b.reader, ByteOrder, &length); err != nil {
		return
	}

	body = make([]byte, length)
	if _, err = b.reader.Read(body); err != nil {
		return
	}

	m = &Message{
		TypeId:        typeId,
		CorrelationId: correlationId,
		Body:          body,
	}
	return
}
