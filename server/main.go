package main

import (
	"flag"
	"github.com/pjvds/publichost"
)

var (
	commandAddress = flag.String("commandAddress", "0.0.0.0:8080", "the command address to serve")
	dataAddress    = flag.String("dataAddress", "0.0.0.0:80", "the data address to serve")
)

func main() {
	flag.Parse()

	log.Info("Starting publichost server")
	log.Info("command address: %v", *commandAddress)
	log.Info("data address: %v", *dataAddress)

	if err := publichost.ListenAndServe(*commandAddress, *dataAddress); err != nil {
		log.Fatal(err)
	}
}
