package protocol

import (
	"fmt"
	"io"
	"strconv"
)

type Request struct {
	Header map[string]string
	Body   io.Reader
}

func NewRequest() *Request {
	return &Request{
		Header: make(map[string]string),
	}
}

// Get the content length in bytes from the request header. If the header
// is not present the a content length of zero is returned and no error.
func getContentLength(request *Request) (int, error) {
	contentLength, ok := request.Header["ContentLength"]
	if !ok {
		return 0, nil
	}

	length, err := strconv.Atoi(contentLength)
	if err != nil || length < 0 {
		return 0, fmt.Errorf("Invalid content length header value")
	}

	return length, nil
}
