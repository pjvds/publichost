package network

import (
	"io"
)

const (
	// Indicator for packet start
	magicStart = byte(0x88)
	// Indicates for packet end
	magicEnd = byte(0x89)
)

type Packet struct {
	data []byte
}

func NewPacket(data []byte) *Packet {
	return &Packet{
		data: data,
	}
}

func (p *Packet) CreateContentReader() io.Reader {
	return bytes.NewBuffer(p.data[1:])
}

func (p *Packet) TypeId() byte {
	return p.data[0]
}

func (p *Packet) Len() int16 {
	return int16(len(p.data))
}

func (p *Packet) WriteTo(writer io.Writer) (int, error) {
	return writer.Write(p.data)
}
