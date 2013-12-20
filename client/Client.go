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

func NewClientConnection(address string) (ClientConnection, error) {
	log.Debug("Connecting to %v", address)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Debug("Connected to %v", address)

	return handshake(conn)
}

// Will handshake with the server. The caller is responsible
// for closing the connection on an error.
func handshake(conn net.Conn) (ClientConnection, error) {
	log.Debug("Starting handshake")

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

	log.Debug("Written handshake request")

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
