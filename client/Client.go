package client

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

type ClientConnection interface {
	GetAddress() string
	Close() error
}

type clientConnection struct {
	publicAddress string // The public address of the tunnel

	connection net.Conn
	reader     *bufio.Reader
	writer     *bufio.Writer
}

func (c *clientConnection) GetAddress() string {
	return c.publicAddress
}

func NewClientConnection(address string) (ClientConnection, error) {
	log.Debug("Connecting to %v", address)

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Debug("Connected to %v", address)

	client, err := handshake(conn)
	if err != nil {
		return nil, err
	}

	return client, nil
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
	log.Debug("Handshake replied with %v", response.Status)

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(response.Status)
	}

	return &clientConnection{
		connection:    conn,
		publicAddress: response.Header.Get("X-Tunnel-Hostname"),
	}, nil
}

func (c *clientConnection) Close() error {
	return c.connection.Close()
}
