package network

type TemporaryError interface {
	Temporary() bool
}

func IsTemporaryError(err error) bool {
	if tempErr, ok := err.(TemporaryError); ok && tempErr.Temporary() {
		return true
	}
	return false
}
