package client

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/pjvds/publichost/message"
	"net"
)

var (
	address       = flag.String("address", "localhost:80", "The local address to make publicly available")
	remoteAddress = "publichost.me:8080"
)

type PublicHostClient struct {
	conn   net.Conn
	reader *message.Reader
	writer *message.Writer
}

func Dial(address string) (*PublicHostClient, error) {
	resolved, err := net.ResolveTCPAddr("tcp", remoteAddress)
	if err != nil {
		return nil, fmt.Errorf("Unable to resolve %v: %v", remoteAddress, err)
	}

	conn, err := net.DialTCP("tcp", nil, resolved)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect publichost server: %v", err)
	}

	bufferedReader := bufio.NewReader(conn)

	return &PublicHostClient{
		conn:   conn,
		reader: message.NewReader(bufferedReader),
	}, nil
}

func (c *PublicHostClient) handshake() error {
	// This is for future use. We don't have
	// anything to exchange during handshake.
	return nil
}

func (p *PublicHostClient) CreateTunnel(hostname string, port int) {
}

func (c *PublicHostClient) request(request message.Message) (response message.Message, err error) {
	if err := c.writer.Write(request); err != nil {
		return nil, err
	}

	response, err = c.reader.Read()
	return
}
