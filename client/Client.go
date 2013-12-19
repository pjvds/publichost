package client

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type ClientConnection interface {
}

type clientConnection struct {
	connection net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
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

// Will handshake with the server. The caller is responsible
// for closing the connection on an error.
func handshake(conn net.Conn) (ClientConnection, error) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	request, _ := http.NewRequest("CONNECT", "tunneler.publichost.me", nil)
	request.Header.Set("Connection", "Keep-Alive")
	request.Header.Set("X-PublicHost-HostName", "firstclient.publichost.me")

	if err := request.Write(writer); err != nil {
		return nil, err
	}
	if err := writer.Flush(); err != nil {
		return nil, err
	}

	response, err := http.ReadResponse(reader, request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Cannot create tunnel: %v (%v)", response.Status, response.StatusCode)
	}

	// TODO: Negociate protocol version, etc.
	return clientConnection{
		connection: conn,
	}, nil
}
