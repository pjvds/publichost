package network

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
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

type Connection interface {
	Send(packet *Packet) (err error)
	Receive() (packet *Packet, err error)
	Close() error
	String() string
}

type connection struct {
	conn net.Conn

	byteOrder binary.ByteOrder

	reader *bufio.Reader
	writer *bufio.Writer
}

// Creates a new publichost connection for a just
// accepted underlying network connection (probably TCP).
func newConnection(conn net.Conn) (Connection, error) {
	return &connection{
		conn:      conn,
		byteOrder: binary.BigEndian,
		reader:    bufio.NewReader(conn),
		writer:    bufio.NewWriter(conn),
	}, nil
}

func Dial(address string) (Connection, error) {
	var conn net.Conn
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	return &connection{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}, nil
}

func (c *connection) String() string {
	return c.conn.RemoteAddr().String()
}

func (c *connection) Send(packet *Packet) (err error) {
	defer func() {
		if err != nil {
			log.Info("error sending packet to %v: %v", c.conn.RemoteAddr(), err)
		}
	}()

	// Write packet to a buffer
	var buffer bytes.Buffer
	if _, err = packet.WriteTo(&buffer); err != nil {
		return
	}

	// Write packet start indicator
	if err = binary.Write(c.writer, c.byteOrder, magicStart); err != nil {
		return
	}

	// Write packet size
	if err = binary.Write(c.writer, c.byteOrder, int16(buffer.Len())); err != nil {
		return
	}

	// Write packet content
	if _, err = buffer.WriteTo(c.conn); err != nil {
		return
	}

	// Write packet end
	err = binary.Write(c.writer, c.byteOrder, magicEnd)
	return
}

func (c *connection) Receive() (packet *Packet, err error) {
	var nextByte byte
	if err = binary.Read(c.reader, c.byteOrder, &nextByte); err != nil {
		return
	}
	if nextByte != magicStart {
		// TODO: Search the next package start
		err = fmt.Errorf("wrong packet start indicator %#x", nextByte)
		return
	}

	var length int16
	if err = binary.Read(c.reader, c.byteOrder, &length); err != nil {
		return
	}

	var data []byte
	bodyReader := io.LimitReader(c.reader, int64(length))
	if data, err = ioutil.ReadAll(bodyReader); err != nil {
		return
	}

	if err = binary.Read(c.reader, c.byteOrder, &nextByte); err != nil {
		return
	}
	if nextByte != magicEnd {
		err = fmt.Errorf("wrong packet end indicator %#x", nextByte)
		return
	}

	packet = NewPacket(data)
	return
}

func (c *connection) Close() error {
	return c.conn.Close()
}
