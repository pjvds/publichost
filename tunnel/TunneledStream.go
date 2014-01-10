package tunnel

import (
	"github.com/pjvds/publichost/stream"
)

type TunneledStream struct {
	streamId stream.Id
	tunnel   Tunnel
}

func NewTunneledStream(streamId stream.Id, tunnel Tunnel) stream.Stream {
	return &TunneledStream{
		streamId: streamId,
		tunnel:   tunnel,
	}
}

func (t *TunneledStream) Read(p []byte) (n int, err error) {
	return t.tunnel.ReadStream(t.streamId, p)
}

func (t *TunneledStream) Write(p []byte) (n int, err error) {
	return t.tunnel.WriteStream(t.streamId, p)
}

func (t *TunneledStream) Close() (err error) {
	return t.tunnel.CloseStream(t.streamId)
}
