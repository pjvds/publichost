package proto

import (
	"encoding/gob"
	"io"
)

type ProxyDataRequest struct {
	StreamId int16

	Flags ProxyDataFlags

	Data []byte
}

func NewProxyDataRequest(streamId int16, flags ProxyDataFlags, data []byte) *ProxyDataRequest {
	return &ProxyDataRequest{
		StreamId: streamId,
		Flags:    flags,
		Data:     data,
	}
}

func (m *ProxyDataRequest) GetTypeId() byte {
	return TypeProxyDataRequest
}

func (m *ProxyDataRequest) Write(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(m)
}
