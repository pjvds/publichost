package tunnel

import (
	"github.com/op/go-logging"
	"github.com/pjvds/publichost/net/message"
	"github.com/pjvds/publichost/stream"
	"bytes"
	"encoding/binary"
)

type ReadStreamHandler struct {
	tunnel Tunnel
	log    *logging.Logger
}

func NewReadStreamHandler(tunnel Tunnel) MessageHandler {
	return &ReadStreamHandler{
		tunnel: tunnel,
		log:    logging.MustGetLogger("handlers"),
	}
}

func (h *ReadStreamHandler) Handle(response ResponseWriter, m *message.Message) error {
	log.Debug("handling read stream request: %v", m.String())
	
	buffer := bytes.NewBuffer(m.Body)
	streamId, err := stream.ReadId(buffer)
	log.Debug("stream id: %v", streamId)

	var size uint32
	binary.Read(buffer, message.ByteOrder, &size)
	p := make([]byte, size)
	n, err := h.tunnel.ReadStream(streamId, p)

	log.Debug("read %v bytes", n)

	if err != nil {
		return response.Nack(err)
	}

	return response.Ack(p[0:n])
}
