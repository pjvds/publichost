package proto

import (
	"encoding/gob"
	"io"
)

const (
	TypeExposeRequest     = byte(iota)
	TypeExposeResponse    = byte(iota)
	TypeProxyDataRequest  = byte(iota)
	TypeProxyDataResponse = byte(iota)
	TypeNokResponse       = byte(iota)
)

type ExposeRequest struct {
	LocalAddress string
}

func NewExposeRequest(localAddress string) *ExposeRequest {
	return &ExposeRequest{
		LocalAddress: localAddress,
	}
}

func (m *ExposeRequest) GetTypeId() byte {
	return TypeExposeRequest
}

func (m *ExposeRequest) Write(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(m)
}

type ExposeResponse struct {
	RouteId       int32
	RemoteAddress string
}

func (m *ExposeResponse) GetTypeId() byte {
	return TypeExposeResponse
}

func (m *ExposeResponse) Write(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(m)
}

type NokResponse struct {
	Error string
}

func (m *NokResponse) GetTypeId() byte {
	return TypeNokResponse
}

func (m *NokResponse) Write(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(m)
}
