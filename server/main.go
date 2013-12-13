package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	address = flag.String("address", "0.0.0.0:8080", "The address to host on")
)

type Tunnel struct {
	Name string
}

type Packet struct {
}

func main() {
	flag.Parse()
	accepted := make(chan *net.TCPConn, 50)

	addr, err := net.ResolveTCPAddr("tcp", *address)
	if err != nil {
		fmt.Printf("Cannot resolve address %v: %v\n", *address, err)
		return
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Printf("Cannot start listening on address %v: %v\n", *address, err)
		return
	}

	defer listener.Close()
	fmt.Printf("Starting to accept connection on %v\n", listener.Addr())

	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		fmt.Printf("Accepted connection from %v", connection.RemoteAddr())
		accepted <- connection
	}
}
