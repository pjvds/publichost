package tunnel

import (
	"github.com/op/go-logging"
	"github.com/pjvds/publichost/net/message"
	"github.com/pjvds/publichost/stream"
)

type CloseStreamHandler struct {
	tunnel Tunnel
	log    *logging.Logger
}

func NewCloseStreamHandler(tunnel Tunnel) MessageHandler {
	return &CloseStreamHandler{
		tunnel: tunnel,
		log:    logging.MustGetLogger("handers"),
	}
}

func (h *CloseStreamHandler) Handle(response ResponseWriter, m *message.Message) error {
	streamId := stream.ParseId(m.Body)
	err := h.tunnel.CloseStream(streamId)

	if err != nil {
		return response.Nack(err)
	}

	return response.Ack(streamId.Bytes())
}
