package publichost

import (
	"bufio"
)

// A connection can receive request- and
// response messages.
type connectionReader struct {
	Requests  chan *Request
	Responses chan *Response

	reader *bufio.Reader
}

func newConnectionReader(reader *bufio.Reader) *ConnectionReader {
	return &ConnectionReader{
		Requests:  make(chan *Request, 25),
		Responses: make(chan *Response, 25),

		reader: reader,
	}
}

func (r *MessageReader) nextMessageType() (messageType byte, err error) {
	messageType, err = r.reader.ReadByte()
	if err != nil {
		if neter, ok := err.(net.Error); ok && neter.Temporary() {
			toSleep := 1 * time.Second

			log.Warning("temporary error reading message type from connection: %v; retry in %v", neter, toSleep)
			time.Sleep()
			continue
		}
		return
	}

	return
}

func (r *MessageReader) Read() error {
	for {
		switch messageType {
		case TRequest:
			request, err := readNextRequest(r.reader)
			if err != nil {
				return err
			}
			r.Requests <- request
		case TResponse:
			response, err := readNextResponse(r.reader)
			if err != nil {
				return err
			}
			r.Responses <- response
		default:
			return fmt.Errorf("unexpected message type flag: %v", messageType)
		}

		request, err := readNextRequest(r.reader)
		if err != nil {
			return err
		}

		r.Requests <- request
	}

	return nil
}
