package main

import (
	"github.com/pjvds/publichost/server"
)

func main() {
	address := ":http"

	log.Info("Starting publichost server at %v", address)

	server := server.NewServer(address)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error serving at %v: %v", address, err)
	}
}
