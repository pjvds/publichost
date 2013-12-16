package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Read a request from the reader.
func ReadRequest(reader io.Reader) (*Request, error) {
	request := &Request{
		Header: make(map[string]string),
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		if isEndOfHeader(line) {
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
	request.Body = io.LimitReader(reader, int64(contentLength))

	return request, nil
}

// Determine if the given line marks the end of the header section of a request.
func isEndOfHeader(line string) bool {
	return len(line) == 0
}
