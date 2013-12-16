package protocol

type RequestHeader map[string]string

// Gets the value, if exists; otherwise, emtpy string.
func (r RequestHeader) Get(name string) string {
	value, ok := r[name]
	if ok {
		return value
	}
	return ""
}

func NewRequestHeader() RequestHeader {
	return make(RequestHeader)
}
