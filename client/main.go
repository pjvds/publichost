package main

import (
	"flag"
	"fmt"
	"github.com/pjvds/publichost/protocol"
	"net"
)

var (
	address       = flag.String("address", "localhost:80", "The local address to make publicly available")
	remoteAddress = "publichost.me:8080"
)

func main() {
	flag.Parse()

	resolved, err := net.ResolveTCPAddr("tcp", remoteAddress)
	if err != nil {
		fmt.Printf("Unable to resolve %v: %v\n", remoteAddress, err)
		return
	}

	connection, err := net.DialTCP("tcp", nil, resolved)
	if err != nil {
		fmt.Printf("Unable to connect publichost server: %v\n", err)
	}

	defer connection.Close()
	fmt.Println("Connected to publichost server")

	for {
		request, err := receive(connection)
		if err != nil {
			fmt.Printf("Error receiving message: %v", err)
		} else {
			fmt.Printf("Received: %s", request)
		}
	}
}

func receive(connection net.Conn) (*protocol.Request, error) {
	message, err := protocol.ReadRequest(connection)
	return message, err
}
