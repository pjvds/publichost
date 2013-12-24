package publichost

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
)

var (
	headerSize = 2 + // Id
		1 + // Type
		2 // Length
)

type Request struct {
	// The identifier of the request. This number
	// should be pseudo unique per tunnel. It gets
	// recycled over time.
	Id uint16

	// Identifies the type of the request.
	Type uint8

	// The length of the body.
	Length uint16

	Body io.Reader
}

func NewRequest(id uint16, t uint8, body io.Reader) *Request {
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

	return &Request{
		Id:     id,
		Type:   t,
		Length: length,
		Body:   body,
	}
}

func (r Request) Write(writer io.Writer) (err error) {
	if err = binary.Write(writer, binary.BigEndian, r.Id); err != nil {
		return
	}
	if err = binary.Write(writer, binary.BigEndian, r.Type); err != nil {
		return
	}
	if err = binary.Write(writer, binary.BigEndian, r.Length); err != nil {
		return
	}

	if r.Body != nil {
		io.Copy(dst, src)
	}
	return
}

func ReadRequest(reader io.Reader) (request *Request, err error) {
	r := new(Request)
	if err = binary.Read(reader, binary.BigEndian, &r.Id); err != nil {
		return
	}
	if err = binary.Read(reader, binary.BigEndian, &r.Type); err != nil {
		return
	}
	if err = binary.Read(reader, binary.BigEndian, &r.Length); err != nil {
		return
	}

	if r.Length > 0 {
		r.Body = io.LimitReader(reader, int64(r.Length))
	}

	request = r
	return
}
