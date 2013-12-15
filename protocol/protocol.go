package protocol

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Message struct {
	Header  map[string]string
	Payload io.Reader
}

func NewMessage() *Message {
	return &Message{
		Header: make(map[string]string),
	}
}

// Read a message from the reader.
func ReadMessage(reader io.Reader) (*Message, error) {
	message := &Message{
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

		message.Header[key] = value
	}

	// Check if the scanner didn't had problems
	// reading from the reader.
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Get content length from the message header.
	contentLength, err := getContentLength(message)
	if err != nil {
		return nil, err
	}

	// TODO: What if we close this reader?
	message.Payload = io.LimitReader(reader, int64(contentLength))

	return message, nil
}

// Determine if the given line marks the end of the header section of a message.
func IsEndOfHeader(line string) bool {
	return len(line) == 0
}

// Get the content length in bytes from the message header. If the header
// is not present the a content length of zero is returned and no error.
func getContentLength(message *Message) (int, error) {
	contentLength, ok := message.Header["ContentLength"]
	if !ok {
		return 0, nil
	}

	length, err := strconv.Atoi(contentLength)
	if err != nil || length < 0 {
		return 0, fmt.Errorf("Invalid content length header value")
	}

	return length, nil
}
