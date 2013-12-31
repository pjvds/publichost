package stream

import (
	"errors"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Map interface {
	// Adds a stream to the map. It returns
	// an ErrAlreadyExists when there is
	// already an entry with the same id.
	Add(id Id, s Stream) (err error)

	// Retrieve a stream by its id. It returns
	// an ErrNotFound when there is no stream
	// added with that id.
	Get(id Id) (s Stream, err error)

	// Deletes the stream from the map. It returns
	// an ErroNotFound when there is no stream
	// with that id.
	Delete(id Id) (err error)
}
