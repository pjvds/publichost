package publichost

import (
	"bufio"
	"net"
)

type Client struct {
	conn net.Conn

	reader *bufio.Reader
	writer *bufio.Writer
}

func DialAndServe(address string) (err error) {
	var conn net.Conn
	if conn, err = net.Dial("tcp", address); err != nil {
		return
	}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	c := &Client{
		conn:   conn,
		reader: reader,
		writer: writer,
	}

	if err = c.handshake(); err != nil {
		defer conn.Close()
		return
	}

	return c.serve()
}

func (c *Client) handshake() (err error) {
	r := NewRequest(0, OpOpenTunnel, nil)
	return c.sendRequest(r)
}

func (c *Client) serve() error {
	for {

	}
}

func (c *Client) sendRequest(request *Request) error {
	if err := request.Write(c.writer); err != nil {
		return err
	}
	if err := c.writer.Flush(); err != nil {
		return err
	}
	return nil
}
