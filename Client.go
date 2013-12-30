package publichost

import (
	"bufio"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"net"
	"sync"
)

const (
	OpOpenStream  = byte(iota)
	OpStreamData  = byte(iota)
	OpCloseStream = byte(iota)
	Ack           = byte(iota)
	Nack          = byte(iota)
)

var (
	ByteOrder = binary.BigEndian
)

type Message struct {
	TypeId        byte
	CorrelationId uint32
	Body          []byte
}

type Client struct {
	sequence *idSequence

	conn net.Conn

	reader *bufio.Reader
	writer *bufio.Writer

	streams     map[StreamId]*StreamFrondEnd
	streamsLock sync.Mutex
}

func (c *Client) Serve() (err error) {
	var message *Message

	for {
		if message, err = c.receiveNext(); err != nil {
			break
		}

		switch message.TypeId {
		case OpOpenStream:
			go c.handleOpenStream(message)
		case OpStreamData:
			go c.handleStreamData(message)
		case OpCloseStream:
			go c.handleCloseStream(message)
		}
	}

	return
}

func (c *Client) handleOpenStream(m *Message) {
	streamId := StreamId(c.sequence.Next())
	address := string(m.Body)

	stream, err := Dial(streamId, address)
	if err != nil {
		err = c.nack(m.CorrelationId, err)
	} else {
		c.streams[streamId] = stream
		c.ack(m.CorrelationId)

		go stream.Serve()
	}
}

func (c *Client) handleStreamData(m *Message) {
	streamId := StreamId(ByteOrder.Uint32(m.Body[0:4]))
	data := m.Body[4:]

	stream, ok := c.streams[streamId]
	if !ok {
		c.nack(m.CorrelationId, fmt.Errorf("stream not found"))
		return
	}

	ack, nack := stream.StreamData(data)
	select {
	case <-ack:
		c.ack(m.CorrelationId)
	case err := <-nack:
		c.nack(m.CorrelationId, err)
	}
}

func (c *Client) handleCloseStream(m *Message) {
	streamId := StreamId(ByteOrder.Uint32(m.Body[0:4]))

	stream, ok := c.streams[streamId]
	if !ok {
		c.nack(m.CorrelationId, fmt.Errorf("stream not found"))
		return
	}

	stream.Close()
	delete(c.streams, streamId)

	c.ack(m.CorrelationId)
}

func (c *Client) ack(correlationId uint32) error {
	m := &Message{
		TypeId:        Ack,
		CorrelationId: correlationId,
	}
	return c.send(m)
}

func (c *Client) nack(correlationId uint32, err error) error {
	m := &Message{
		TypeId:        Nack,
		CorrelationId: correlationId,
		Body:          []byte(err.Error()),
	}
	return c.send(m)
}

func (c *Client) receiveNext() (m *Message, err error) {
	decoder := gob.NewDecoder(c.reader)
	err = decoder.Decode(m)
	return
}

func (c *Client) send(m *Message) error {
	encoder := gob.NewEncoder(c.writer)
	return encoder.Encode(m)
}
