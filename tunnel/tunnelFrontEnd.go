package tunnel

import (
	"encoding/binary"
	"errors"
	"github.com/pjvds/publichost/net"
	"github.com/pjvds/publichost/net/message"
	"github.com/pjvds/publichost/stream"
)

type tunnelFrontEnd struct {
	conn            net.ClientConnection
	messageSequence message.IdSequence
}

func NewTunnelFrontEnd(conn net.ClientConnection) Tunnel {
	return &tunnelFrontEnd{
		conn:            conn,
		messageSequence: message.NewIdSequence(),
	}
}

func (t *tunnelFrontEnd) OpenStream(network, address string) (id stream.Id, err error) {
	var response *message.Message

	// TODO: use difference sequence
	request := message.NewMessage(message.OpOpenStream, uint64(t.messageSequence.Next()), []byte(address))
	if response, err = t.conn.SendRequest(request); err != nil {
		return
	}

	switch response.TypeId {
	case message.Ack:
		id = stream.Id(binary.BigEndian.Uint32(response.Body))
	case message.Nack:
		err = errors.New(string(response.Body))
	default:
		log.Error("unknown response message type: %v", response.TypeId)
		err = errors.New("protocol error")
	}
	return
}

func (t *tunnelFrontEnd) ReadStream(id stream.Id, p []byte) (n int, err error) {
	var response *message.Message

	// TODO: use difference sequence
	request := message.NewMessage(message.OpReadStream, uint64(t.messageSequence.Next()), []byte(address))
	if response, err = t.conn.SendRequest(request); err != nil {
		return
	}

	switch response.TypeId {
	case message.Ack:
		id = stream.Id(binary.BigEndian.Uint32(response.Body))
	case message.Nack:
		err = errors.New(string(response.Body))
	default:
		log.Error("unknown response message type: %v", response.TypeId)
		err = errors.New("protocol error")
	}
	return
}

func (t *tunnelFrontEnd) WriteStream(id stream.Id, p []byte) (n int, err error) {
	var s stream.Stream

	if s, err = t.streams.Get(id); err != nil {
		return
	}

	// TODO: Should be remove the stream on an error?
	return s.Write(p)
}

func (t *tunnelFrontEnd) CloseStream(id stream.Id) (err error) {
	var s stream.Stream

	if s, err = t.streams.Get(id); err != nil {
		return
	}

	if err = s.Close(); err != nil {
		return
	}

	t.streams.Delete(id)
	return
}
