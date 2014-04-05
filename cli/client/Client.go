package main

import (
	"flag"
	"net"
	"github.com/pjvds/publichost/tunnel"
)

var (
	address = flag.String("address", "localhost:8081", "The publichost server to connect")

	localAddress = flag.String("local", "localhost:4000", "The local address to expose")
)

func main() {
	flag.Parse()

	conn, err := net.Dial("tcp4", *address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	host := tunnel.NewBackendHost(conn, *localAddress)
	if err := host.Serve(); err != nil {
		log.Fatal(err)
	}
}
