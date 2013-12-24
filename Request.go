package publichost

import (
	"encoding/binary"
	"io"
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
