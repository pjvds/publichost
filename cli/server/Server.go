package main

import (
	"flag"
	"github.com/pjvds/publichost/server"
)

var (
	address = flag.String("address", "127.0.0.1:8080", "The address to bind to")
)

func main() {
	flag.Parse()

	println("publichost - v0.1")
	println("hosting at: " + *address)

	if err := server.ListenAndServe(*address); err != nil {
		println("err: " + err.Error())
	}
}
