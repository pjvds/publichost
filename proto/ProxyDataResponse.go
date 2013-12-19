package proto

import (
	"encoding/gob"
	"io"
)

type ProxyDataResponse struct {
	Status byte
}

func NewProxyDataResponse(status byte) *ProxyDataResponse {
	return &ProxyDataResponse{
		Status: status,
	}
}

func (m *ProxyDataResponse) GetTypeId() byte {
	return TypeProxyDataResponse
}

func (m *ProxyDataResponse) Write(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(m)
}
