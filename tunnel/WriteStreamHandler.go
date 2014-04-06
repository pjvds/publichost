package tunnel

import (
    "github.com/op/go-logging"
    "github.com/pjvds/publichost/net/message"
    "github.com/pjvds/publichost/stream"
)

type WriteStreamHandler struct {
    tunnel Tunnel
    log    *logging.Logger
}

func NewWriteStreamHandler(tunnel Tunnel) MessageHandler {
    return &WriteStreamHandler{
        tunnel: tunnel,
        log:    logging.MustGetLogger("handers"),
    }
}

func (h *WriteStreamHandler) Handle(response ResponseWriter, m *message.Message) error {
    log.Debug("handling write stream request: %v", m)

    streamId := stream.ParseId(m.Body[0:4])
    data := m.Body[4:]

    log.Debug("stream id: %v", streamId)
    log.Debug("data: %v bytes", len(data))

    n, err := h.tunnel.WriteStream(streamId, data)

    if err != nil {
        return response.Nack(err)
    }

    body := make([]byte, 4, 4)
    message.ByteOrder.PutUint32(body, uint32(n))

    return response.Ack(body)
}
