package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Request struct {
	Header map[string]string
	Body   io.Reader
}

// Read a request from the reader.
func ReadRequest(reader io.Reader) (*Request, error) {
	request := &Request{
		Header: make(map[string]string),
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if IsEndOfHeader(line) {
			break
		}

		data := strings.SplitN(line, ":", 2)
		if len(data) != 2 {
			return nil, fmt.Errorf("Invalid header entry: %v", line)
		}

		key := data[0]
		value := data[1]

		request.Header[key] = value
	}

	// Check if the scanner didn't had problems
	// reading from the reader.
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Get content length from the request header.
	contentLength, err := getContentLength(request)
	if err != nil {
		return nil, err
	}

	// TODO: What if we close this reader?
	request.Body = io.LimitReader(reader, contentLength)

	return request, nil
}

// Determine if the given line marks the end of the header section of a request.
func IsEndOfHeader(line string) bool {
	return len(line) == 0
}

// Get the content length in bytes from the request header. If the header
// is not present the a content length of zero is returned and no error.
func getContentLength(request *Request) (int64, error) {
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

type TunnelResponse struct {
	Headers        map[string]string
	StatusCode     int
	StatusResponse string
	Body           io.Reader
}
