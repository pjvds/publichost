package main

import (
	"flag"
	"github.com/pjvds/publichost/net"
)

var (
	address = flag.String("address", "publichost.me:8080", "The address to bind to")
)

func main() {
	flag.Parse()

	conn, err := net.Dial(*address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	host := tunnel.NewBackendHost(conn)
	if err := host.Serve(); err != nil {
		log.Fatal(err)
	}
}
