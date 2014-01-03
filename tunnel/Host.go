package tunnel

import (
	"errors"
	"github.com/pjvds/publichost/net/message"
	"io"
)

type Host interface {
	Serve() error
}

type host struct {
	reader message.Reader
	writer message.Writer

	tunnel Tunnel

	handlers map[byte]MessageHandler
}

func NewTunnelHost(conn io.ReadWriteCloser) Host {
	h := &host{
		reader:   message.NewReader(conn),
		writer:   message.NewWriter(conn),
		tunnel:   NewTunnelBackend(),
		handlers: make(map[byte]MessageHandler),
	}
	h.handlers[message.OpOpenStream] = NewOpenStreamHandler(h.tunnel)

	return h
}

func (h *host) Serve() (err error) {
	var request *message.Message

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
			go response.Nack(errors.New("unknown type id"))
		}
	}

	return
}
