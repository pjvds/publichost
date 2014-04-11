package message

import (
	"bufio"
	"encoding/binary"
)

type Reader interface {
	Read() (m *Message, err error)
}

type bufferedReader struct {
	reader *bufio.Reader
}

func NewReader(r *bufio.Reader) Reader {
	return &bufferedReader{
		reader: r,
	}
}

func (b *bufferedReader) Read() (m *Message, err error) {
	var firstByte byte
	var typeId byte
	var correlationId uint64
	var length uint16
	var body []byte

	if err = binary.Read(b.reader, ByteOrder, &firstByte); err != nil {
		return
	}

	if firstByte != MagicStart {
		log.Warning("first byte missmatch: got %v, expected %v", firstByte, MagicStart)		
		if err = b.readUntilMessageStart(); err != nil {
			return
		}
	}

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
	log.Debug("message received: %v", m)

	return
}

func (b *bufferedReader) readUntilMessageStart() (err error){
	var firstByte byte
	if firstByte, err = b.reader.ReadByte(); err != nil{
		return
	}

	for firstByte != MagicStart {
		if firstByte, err = b.reader.ReadByte(); err != nil{
			return
		}	
	}

	return
}