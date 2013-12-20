package main

import (
	"flag"
	"fmt"
	"github.com/pjvds/publichost/client"
)

var (
	address = flag.String("address", "tunneler.publichost.me", "the address of the publichost frond end service")
)

func main() {
	flag.Parse()

	_, err := client.NewClientConnection(*address)
	fmt.Printf("Err: %v\n", err)
}
