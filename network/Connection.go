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

type Packet interface {
	Write(writer io.Writer) error
}

type Listener struct {
	listener net.Listener
}

func Listen(address string) (listener *Listener, err error) {
	var netListener net.Listener
	if netListener, err = net.Listen("tcp", address); err != nil {
		return
	}

	listener = &Listener{
		listener: netListener,
	}
	return
}

type Connection struct {
	conn net.Conn

	byteOrder binary.ByteOrder

	writer *bufio.Writer
	reader *bufio.Reader
}

func Dial(address string) (conn *Connection, err error) {
	var underlyingConn net.Conn
	if underlyingConn, err = net.Dial("tcp", address); err != nil {
		return
	}

	conn = &Connection{
		conn:   underlyingConn,
		reader: bufio.NewReader(underlyingConn),
		writer: bufio.NewWriter(underlyingConn),
	}
	return
}

func (c *Connection) Send(packet Packet) (err error) {
	defer func() {
		if err != nil {
			log.Info("error sending packet to %v: %v", c.conn.RemoteAddr(), err)
		}
	}()

	// Write packet to a buffer
	var buffer bytes.Buffer
	if err = packet.Write(&buffer); err != nil {
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

func (c *Connection) Receive() (packet *Packet, err error) {
	var nextByte byte
	if err = c.reader.Read(reader, c.byteOrder, &nextByte); err != nil {
		return
	}
	if nextByte != magicStart {
		// TODO: Search the next package start
		return fmt.Errorf("wrong packet start indicator %#x", nextByte)
	}

	var length int16
	if err = binary.Read(reader, c.byteOrder, data); err != nil {
		return
	}

	var data []byte
	reader := io.LimitReader(reader, length)
	if data, err = ioutil.ReadAll(reader); err != nil {
		return
	}

	if err = c.reader.Read(reader, c.byteOrder, &nextByte); err != nil {
		return
	}
	if nextByte != magicEnd {
		return fmt.Errorf("wrong packet end indicator %#x", nextByte)
	}
}

func (c *Connection) Close() error {
	return c.conn.Close()
}
