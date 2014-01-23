package tunnel

import (
	"github.com/pjvds/publichost/stream"
)

type backend struct {
	dialer    stream.Dialer
	sequencer stream.IdSequence
	streams   stream.Map
}

func NewTunnelBackend() Tunnel {
	return &backend{
		dialer:    stream.NewDialer(),
		sequencer: stream.NewIdSequence(),
		streams:   stream.NewThreadSafeMap(),
	}
}

func (t *backend) OpenStream(network, address string) (id stream.Id, err error) {
	var s stream.Stream
	if s, err = t.dialer.Dial(network, address); err != nil {
		return
	}

	id = t.sequencer.Next()
	t.streams.Add(id, s)

	return
}

func (t *backend) ReadStream(id stream.Id, p []byte) (n int, err error) {
	var s stream.Stream

	if s, err = t.streams.Get(id); err != nil {
		return
	}

	// TODO: Should be remove the stream on an error?
	return s.Read(p)
}

func (t *backend) WriteStream(id stream.Id, p []byte) (n int, err error) {
	var s stream.Stream

	if s, err = t.streams.Get(id); err != nil {
		return
	}

	// TODO: Should be remove the stream on an error?
	return s.Write(p)
}

func (t *backend) CloseStream(id stream.Id) (err error) {
	var s stream.Stream

	if s, err = t.streams.Get(id); err != nil {
		return
	}

	if err = s.Close(); err != nil {
		return
	}

	t.streams.Delete(id)
	return
}
