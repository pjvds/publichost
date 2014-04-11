package net

import (
	"github.com/pjvds/publichost/net/message"
	"io"
	"net"
	"sync"
	"time"
	"errors"
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

func Dial(address string, timeout time.Duration) (c ClientConnection, err error) {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return nil, err
	}

	c = NewClientConnection(conn)
	return
}

type timeoutWrapper struct{
	conn net.Conn
	timeout time.Duration
}

func (t *timeoutWrapper) Write(p []byte) (n int, err error) {
	// absolute time after which I/O operations fail with a timeout instead of blocking
	if err = t.conn.SetWriteDeadline(time.Now().Add(t.timeout)); err != nil {
		log.Error("error setting read deadline on connection: %v", err)
		return
	}
	n, err = t.conn.Write(p)
	return
}

func NewClientConnection(conn net.Conn) (c ClientConnection) {
	timeoutWrapper := &timeoutWrapper{
		timeout: time.Second,
		conn: conn,
	}

	c = &clientConnection{
		conn:   conn,
		reader: message.NewReader(conn),
		writer: message.NewWriter(timeoutWrapper),

		wg: &sync.WaitGroup{},

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
				log.Warning("error writing request: %v", err)
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
			if neterr, ok := err.(net.Error); ok && neterr.Temporary() {
				log.Warning("temporary error reading connection: %v", err)
				continue
			}
			log.Fatalf("error reading connection: %v", err)
			return
		}

		if r, ok := c.outstandingRequests[response.CorrelationId]; ok {
			r.Response <- response
		} else {
			log.Warning("received message for unknown request: %v", response.CorrelationId)
		}
	}
}

func (c *clientConnection) SendRequest(request *message.Message) (response *message.Message, err error) {
	r := &roundtrip{
		request:  request,
		Response: make(chan *message.Message, 1),
		Error:    make(chan error),
	}
	c.outstandingRequests[request.CorrelationId] = r
	defer delete(c.outstandingRequests, request.CorrelationId)

	c.outgoing <- r
	log.Debug("send request %v", request.String())

	select {
	case response = <-r.Response:
		log.Debug("received response: %v", response.String())
		return
	case err = <-r.Error:
		log.Debug("error: %v", response.String())
		return
	case <-time.After(5 * time.Second):
		log.Warning("timeout: %v", request)
		err = errors.New("timeout")
		return
	}
}
