package main

import (
	"flag"
	"github.com/pjvds/publichost/tunnel"
	"fmt"
)

var (
	address = flag.String("address", "localhost:8080", "The publichost server to connect")
	localAddress = flag.String("local", "localhost:4000", "The local address to expose")
)

func main() {
	flag.Parse()
	fmt.Printf("publichost - v0.1\n")
	fmt.Printf("local address: %v\n", *localAddress)

	host, err := tunnel.NewBackendHost(*address)
	if err != nil {
		log.Fatal(err)
	}

	hostname, err := host.OpenTunnel(*localAddress)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Notice("tunnel opened at: %v", hostname)

	if err := host.Serve(); err != nil {
		log.Fatal(err)
		return
	}
}