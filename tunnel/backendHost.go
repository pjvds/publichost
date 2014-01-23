package tunnel

import (
	"errors"
	"github.com/pjvds/publichost/net/message"
	"io"
)

type backendHost struct {
	conn   io.ReadWriteCloser
	reader message.Reader
	writer message.Writer

	tunnel       Tunnel
	localAddress string

	handlers map[byte]MessageHandler
}

func NewBackendHost(conn io.ReadWriteCloser) Host {
	h := &backendHost{
		conn:     conn,
		reader:   message.NewReader(conn),
		writer:   message.NewWriter(conn),
		tunnel:   NewTunnelBackend(),
		handlers: make(map[byte]MessageHandler),
	}
	h.handlers[message.OpOpenStream] = NewOpenStreamHandler(h.tunnel)
	h.handlers[message.OpCloseStream] = NewCloseStreamHandler(h.tunnel)

	return h
}

func (h *backendHost) Serve() (err error) {
	defer h.conn.Close()

	var request *message.Message
	var response *message.Message

	openTunnel := message.NewMessage(message.OpOpenTunnel, 1, []byte(h.localAddress))
	if err = h.writer.Write(openTunnel); err != nil {
		log.Debug("unable to write handshake message: %v", err)
		return
	}

	if response, err = h.reader.Read(); err != nil {
		log.Debug("unable to read response to handshake: %v", err)
		return
	}

	if response.TypeId == message.Nack {
		err = errors.New(string(response.Body))
		log.Debug("unable to open tunnel: %v", err)
		return
	}

	log.Debug("tunnel opened successfully")

	for {
		if request, err = h.reader.Read(); err != nil {
			if err != io.EOF {
				err = nil
			}
			break
		}

		response := NewResponseWriter(h.writer, request.CorrelationId)

		if handler, ok := h.handlers[request.TypeId]; ok {
			go handler.Handle(response, request)
		} else {
			log.Debug("No handlers for message type id %v", request.TypeId)
			go response.Nack(errors.New("unknown type id"))
		}
	}

	return
}
