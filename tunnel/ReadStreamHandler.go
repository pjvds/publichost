package tunnel

import (
	"bytes"
	"github.com/op/go-logging"
	"github.com/pjvds/publichost/net/message"
	"github.com/pjvds/publichost/stream"
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
	streamId := stream.ParseId(m.Body)

	p := make([]byte, 4098)
	n, err := h.tunnel.ReadStream(streamId, p)

	if err != nil {
		return response.Nack(err)
	}

	return response.Ack(p[0:n])
}
