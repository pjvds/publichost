package stream

import (
	"io"
	"net"
)

type Stream interface {
	io.ReadWriteCloser
}

func Dial(network, address string) (s Stream, err error) {
	return net.Dial(network, address)
}

// // Proxies read requests
// // to a remote reader
// type RemoteReader struct {
// 	StreamId Id

// 	eof bool
// }

// func (r *RemoteReader) IsEOF() bool {
// 	return r.eof
// }

// func (r *RemoteReader) Read(p []byte) (n int, err error) {
// 	if r.IsEOF() {
// 		err = io.EOF
// 		return
// 	}

// 	var response *ReadStreamResponse
// 	request = r.requests.NewReadStreamRequest(id, len(p))

// 	if response, err = request.Execute(); err != nil {
// 		return
// 	}

// 	// Copy response data to p
// 	for i, b := range response.Data {
// 		p[i] = b
// 	}
// 	n = response.Length

// 	if response.EOF {
// 		err = io.EOF
// 	}

// 	return len()
// }

// Returns true if the error indicates a timeout happened.
func IsTimeout(err error) bool {
	switch err := err.(type) {
	case net.Error:
		return err.Timeout()
	default:
		return false
	}
}
