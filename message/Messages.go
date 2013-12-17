package message

import (
	"encoding/gob"
	"io"
)

const (
	TypeExposeRequest    = iota
	TypeExposeResponse   = iota
	TypeProxyDataRequest = iota
	TypeNokResponse      = iota
)

type ExposeRequest struct {
	LocalAddress string
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
