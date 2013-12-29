package server

import (
	"io"
)

type StreamConnection interface {
	io.ReadWriteCloser
}
