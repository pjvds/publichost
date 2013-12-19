package client

import (
	"fmt"
	"net"
	"net/http"
)

type ClientConnection interface {
}

type clientConnection struct {
	connection net.Conn
}

func NewClientConnection(hostname string) (ClientConnection, error) {
	address, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%v:80", hostname))
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, address)
	if err != nil {
		return nil, err
	}

	return handshake(conn)
}

func handshake(conn net.Conn) (ClientConnection, error) {
	request, _ := http.NewRequest("CONNECT", "tunneler.publichost.me", nil)
	request.Header.Set("hostname", "firstclient.publichost.me")

	request.Write(conn)

	// TODO: Negociate protocol version, etc.
	return clientConnection{
		connection: conn,
	}, nil
}
