	package tunnel

import (
	"encoding/binary"
	"errors"
	"github.com/pjvds/publichost/net"
	"github.com/pjvds/publichost/net/message"
	"github.com/pjvds/publichost/stream"
	"bytes"
)

type frondend struct {
	conn            net.ClientConnection
	messageSequence message.IdSequence
}

func NewFrondend(conn net.ClientConnection) Tunnel {
	return &frondend{
		conn:            conn,
		messageSequence: message.NewIdSequence(),
	}
}

func (t *frondend) OpenStream(network, address string) (id stream.Id, err error) {
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

func (t *frondend) ReadStream(id stream.Id, p []byte) (n int, err error) {
	var response *message.Message

	buffer := bytes.NewBuffer(id.Bytes())
	binary.Write(buffer, message.ByteOrder, uint32(len(p)))

	// TODO: use difference sequence
	request := message.NewMessage(message.OpReadStream, uint64(t.messageSequence.Next()), buffer.Bytes())
	if response, err = t.conn.SendRequest(request); err != nil {
		return
	}

	switch response.TypeId {
	case message.Ack:
		id = stream.Id(binary.BigEndian.Uint32(response.Body))
		n = len(response.Body)

		for k, v := range response.Body {
			p[k]=v
		}

	case message.Nack:
		err = errors.New(string(response.Body))
	default:
		log.Error("unknown response message type: %v", response.TypeId)
		err = errors.New("protocol error")
	}
	return
}

func (t *frondend) WriteStream(id stream.Id, p []byte) (n int, err error) {
	var response *message.Message

	// TODO: use difference sequence
	var body bytes.Buffer
	body.Write(id.Bytes())
	body.Write(p)

	request := message.NewMessage(message.OpWriteStream, uint64(t.messageSequence.Next()), body.Bytes())
	if response, err = t.conn.SendRequest(request); err != nil {
		return
	}

	switch response.TypeId {
	case message.Ack:
		p = response.Body
	case message.Nack:
		err = errors.New(string(response.Body))
	default:
		log.Error("unknown response message type: %v", response.TypeId)
		err = errors.New("protocol error")
	}
	return
}

func (t *frondend) CloseStream(id stream.Id) (err error) {
	var response *message.Message

	// TODO: use difference sequence
	request := message.NewMessage(message.OpCloseStream, uint64(t.messageSequence.Next()), id.Bytes())
	if response, err = t.conn.SendRequest(request); err != nil {
		return
	}

	switch response.TypeId {
	case message.Ack:
		// Nothing
		break
	case message.Nack:
		err = errors.New(string(response.Body))
	default:
		log.Error("unknown response message type: %v", response.TypeId)
		err = errors.New("protocol error")
	}
	return
}
