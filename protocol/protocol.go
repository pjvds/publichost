package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Request struct {
	Header map[string]string
	Body   io.Reader
}

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
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return request, nil
}

func IsEndOfHeader(line string) bool {
	return len(line) == 0
}

type TunnelResponse struct {
	Headers        map[string]string
	StatusCode     int
	StatusResponse string
	Body           io.Reader
}
