package main

import (
	"fmt"
	"github.com/pjvds/publichost/client"
)

func main() {
	_, err := client.NewClientConnection("publichost.me")

	fmt.Printf("Err: %v\n", err)
}
