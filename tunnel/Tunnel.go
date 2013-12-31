package tunnel

import (
	"github.com/pjvds/publichost/stream"
)

type Tunnel interface {
	// Opens a new stream at the other end of the tunnel.
	// It returns the id of the stream if created successfully;
	// otherwise, an error is returned.
	OpenStream(network, address string) (id stream.Id, err error)

	ReadStream(id stream.Id, p []byte) (n int, err error)

	WriteStream(id stream.Id, p []byte) (n int, err error)

	CloseStream(id stream.Id) (err error)
}
