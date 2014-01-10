package main

import (
	"flag"
	"github.com/pjvds/publichost/tunnel"
	"log"
	"net"
)

var (
	address = flag.String("address", "publichost.me:8080", "The address to bind to")
)

func main() {
	flag.Parse()

	conn, err := net.Dial("tcp", *address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	host := tunnel.NewTunnelHost(conn)
	if err := host.Serve(); err != nil {
		log.Fatal(err)
	}
}
