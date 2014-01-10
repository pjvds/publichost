package net

import (
	"github.com/pjvds/publichost/net/message"
	"io"
	"net"
	"sync"
)

type ClientConnection interface {
	SendRequest(request *message.Message) (response *message.Message, err error)
	Close() error
}

type roundtrip struct {
	request  *message.Message
	Response chan *message.Message
	Error    chan error
}

type clientConnection struct {
	// This channel will be signaled when the connection is closing.
	closing chan bool

	wg *sync.WaitGroup

	conn   io.ReadWriteCloser
	reader message.Reader
	writer message.Writer

	outstandingRequests map[uint64]*roundtrip

	outgoing chan *roundtrip
}

func Dial(address string) (c ClientConnection, err error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	c = &clientConnection{
		conn:   conn,
		reader: message.NewReader(conn),
		writer: message.NewWriter(conn),

		outstandingRequests: make(map[uint64]*roundtrip),
		outgoing:            make(chan *roundtrip, 10),
	}

	go c.(*clientConnection).serveIncomming()
	go c.(*clientConnection).serveOutgoing()

	return
}

func (c *clientConnection) Close() error {
	return c.conn.Close()
}

func (c *clientConnection) serveOutgoing() error {
	defer c.wg.Done()
	c.wg.Add(1)

	for {
		r := <-c.outgoing
		if r != nil {
			if err := c.writer.Write(r.request); err != nil {
				r.Error <- err
				continue
			}
		}
	}
}

func (c *clientConnection) serveIncomming() (err error) {
	for {
		var response *message.Message
		if response, err = c.reader.Read(); err != nil {
			return
		}

		if r, ok := c.outstandingRequests[response.CorrelationId]; ok {
			r.Response <- response
		}
	}
}

func (c *clientConnection) SendRequest(request *message.Message) (response *message.Message, err error) {
	r := &roundtrip{
		request:  request,
		Response: make(chan *message.Message, 1),
		Error:    make(chan error),
	}
	c.outgoing <- r

	select {
	case response = <-r.Response:
		return
	case err = <-r.Error:
		return
	}
}
