package network

import (
	"net"
)

type Listener interface {
	String() string
	Accept() (Connection, error)
	Close() error
}

type listener struct {
	listener net.Listener
}

func Listen(address string) (l Listener, err error) {
	var netListener net.Listener
	if netListener, err = net.Listen("tcp", address); err != nil {
		return
	}

	l = &listener{
		listener: netListener,
	}
	return
}

func (l *listener) String() string {
	return l.listener.Addr().String()
}

func (l *listener) Accept() (Connection, error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	return newConnection(conn)
}

func (l *listener) Close() error {
	return l.listener.Close()
}
