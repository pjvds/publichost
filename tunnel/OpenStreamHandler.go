package tunnel

import (
	"github.com/op/go-logging"
	"github.com/pjvds/publichost/net/message"
)

type OpenStreamHandler struct {
	localAddress string
	tunnel Tunnel
	log    *logging.Logger
}

func NewOpenStreamHandler(tunnel Tunnel, localAddress string) MessageHandler {
	return &OpenStreamHandler{
		tunnel: tunnel,
		localAddress: localAddress,
		log:    logging.MustGetLogger("handlers"),
	}
}

func (h *OpenStreamHandler) Handle(response ResponseWriter, m *message.Message) error {
	streamId, err := h.tunnel.OpenStream("tcp", h.localAddress)

	if err != nil {
		return response.Nack(err)
	}

	return response.Ack(streamId.Bytes())
}
