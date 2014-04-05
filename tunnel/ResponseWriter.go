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

func (r *responseWriter) Ack(body []byte) (err error) {
	log.Debug("acking to correlation id: %v", r.correlationId)
	m := message.NewMessage(message.Ack, r.correlationId, body)
	
	if err = r.writer.Write(m); err != nil {
		log.Debug("unable to write ack response: %v", err)
	}

	return
}

func (r *responseWriter) Nack(err error) error {
	log.Debug("nacking to correlation id %v: %v", r.correlationId, err)
	m := message.NewMessage(message.Nack, r.correlationId, []byte(err.Error()))
	return r.writer.Write(m)
}
