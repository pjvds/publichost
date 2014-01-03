package tunnel

import (
	"github.com/pjvds/publichost/net/message"
)

type ResponseWriter interface {
	Ack(body []byte) error
	Nack(err error) error
}

type responseWriter struct {
	writer        message.Writer
	correlationId uint64
}

func NewResponseWriter(writer message.Writer, correlationId uint64) ResponseWriter {
	return &responseWriter{
		writer:        writer,
		correlationId: correlationId,
	}
}

func (r *responseWriter) Ack(body []byte) error {
	m := message.NewMessage(message.Ack, r.correlationId, body)
	return r.writer.Write(m)
}

func (r *responseWriter) Nack(err error) error {
	m := message.NewMessage(message.Nack, r.correlationId, []byte(err.Error()))
	return r.writer.Write(m)
}
