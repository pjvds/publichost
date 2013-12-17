package message

import (
	"encoding/gob"
	"io"
)

type Reader struct {
	reader  io.Reader
	decoder Decoder
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		reader:  r,
		decoder: NewDecoder(),
	}
}

func (r *Reader) Read() (Message, error) {
	var header Header

	decoder := gob.NewDecoder(r.reader)
	if err := decoder.Decode(&header); err != nil {
		return nil, err
	}

	result, err := r.decoder.Decode(header.TypeId, header.Length, io.LimitReader(r.reader, int64(header.Length)))
	return result, err
}
