package publichost

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
)

type Response struct {
	// The identifier of the request. This number
	// should be pseudo unique per tunnel. It gets
	// recycled over time.
	RequestId uint16

	StatusCode byte

	// The length of the body.
	Length uint16

	Body io.Reader
}

func NewResponse(requestId uint16, statusCode byte, body io.Reader) *Response {
	length := uint16(0)

	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			length = uint16(v.Len())
		case *bytes.Reader:
			length = uint16(v.Len())
		case *strings.Reader:
			length = uint16(v.Len())
		default:
			panic("body type not supported")
		}
	}

	return &Response{
		RequestId:  requestId,
		StatusCode: statusCode,
		Length:     length,
		Body:       body,
	}
}

func (r *Response) Write(writer io.Writer) (err error) {
	if err = binary.Write(writer, binary.BigEndian, r.RequestId); err != nil {
		return
	}
	if err = binary.Write(writer, binary.BigEndian, r.StatusCode); err != nil {
		return
	}
	if err = binary.Write(writer, binary.BigEndian, r.Length); err != nil {
		return
	}

	if r.Body != nil {
		if _, err = io.Copy(writer, r.Body); err != nil {
			return
		}
	}

	return
}
