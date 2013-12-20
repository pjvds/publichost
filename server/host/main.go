package main

import (
	"flag"
	"github.com/pjvds/publichost/server"
)

var (
	address = flag.String("address", "0.0.0.0:80", "the address to listen on")
)

func main() {
	flag.Parse()
	log.Info("Starting publichost server at %v", *address)

	server := server.NewServer(*address)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error serving at %v: %v", *address, err)
	}
}
