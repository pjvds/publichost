package main

import (
	"flag"
	"fmt"
	"net"
	"time"
)

var (
	address       = flag.String("address", "localhost:80", "The local address to make publicly available")
	remoteAddress = "publichost.me:8080"
)

func main() {
	flag.Parse()

	resolved, err := net.ResolveTCPAddr("tcp", remoteAddress)
	if err != nil {
		fmt.Printf("Unable to resolve %v: %v\n", remoteAddress, err)
		return
	}

	connection, err := net.DialTCP("tcp", nil, resolved)
	if err != nil {
		fmt.Printf("Unable to connect publichost server: %v\n", err)
	}

	defer connection.Close()
	fmt.Println("Connected to publichost server")

	buffer := []byte("CREATE")
	_, err = connection.Write(buffer)
	if err != nil {
		fmt.Println("Error writing data: %v", err)
	}

	time.Sleep(time.Second * 5)

	buffer = []byte("CREATE2")
	_, err = connection.Write(buffer)
	if err != nil {
		fmt.Println("Error writing data: %v", err)
	}
}
