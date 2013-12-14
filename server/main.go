package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
)

var (
	hostname = flag.String("hostname", "publichost.me", "The domain to host on")
)

type Tunnel struct {
	host       *TunnelServiceHost
	connection net.Conn
}

func NewTunnel(host *TunnelServiceHost, connection net.Conn) *Tunnel {
	return &Tunnel{
		host:       host,
		connection: connection,
	}
}

func (t *Tunnel) HandleHttp(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusTeapot)
}

func (t *Tunnel) Serve() {
	buffer := make([]byte, 4096)
	defer t.connection.Close()

	for {
		n, err := t.connection.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading from %v: %v\n", t.connection.RemoteAddr(), err)
			break
		}

		msg := buffer[0:n]
		fmt.Printf("Recieved data from %v: %v\n", t.connection.RemoteAddr(), string(msg))
	}

	fmt.Printf("Closing tunnel from %v")
}

type TunnelServiceHost struct {
}

func main() {
	host := &TunnelServiceHost{}
	if err := host.ListenAndServePH(); err != nil {
		fmt.Sprintf("[FATAL] Error serving PH: %v", err)
	}
}

func (t *TunnelServiceHost) ListenAndServePH() error {
	address := "0.0.0.0:8080"
	resolved, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return fmt.Errorf("Cannot resolve address %v: %v\n", address, err)
	}

	listener, err := net.ListenTCP("tcp", resolved)
	if err != nil {
		return fmt.Errorf("Cannot start listening on address %v: %v\n", resolved, err)
	}

	defer listener.Close()
	fmt.Printf("Starting to accept connection on %v\n", listener.Addr())

	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		fmt.Printf("Accepted connection from %v to %v:%v\n", connection.RemoteAddr(), connection.LocalAddr().Network(), connection.LocalAddr().String())

		tunnel := NewTunnel(t, connection)
		go tunnel.Serve()
	}
}

// func (t *TunnelServiceHost) ListenAndServeHTTP() error {
// 	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
// 		hostname := request.URL.Host

// 		t.rw.Lock()
// 		defer t.rw.Unlock()

// 		if tunnel, ok := t.tunnels[hostname]; ok {
// 			tunnel.HandleHttp(response, request)
// 		} else {
// 			response.WriteHeader(http.StatusNotFound)
// 		}
// 	})
// 	return http.ListenAndServe("", nil)
// }
