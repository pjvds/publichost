package tunnel

import (
	"github.com/op/go-logging"
	"github.com/pjvds/publichost/net/message"
)

type OpenStreamHandler struct {
	tunnel Tunnel
	log    *logging.Logger
}

func NewOpenStreamHandler(tunnel Tunnel) MessageHandler {
	return &OpenStreamHandler{
		tunnel: tunnel,
		log:    logging.MustGetLogger("handlers"),
	}
}

func (h *OpenStreamHandler) Handle(response ResponseWriter, m *message.Message) error {
	address := string(m.Body)
	streamId, err := h.tunnel.OpenStream("tcp", address)

	if err != nil {
		return response.Nack(err)
	}

	return response.Ack(streamId.Bytes())
}
