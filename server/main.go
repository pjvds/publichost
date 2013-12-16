package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/pjvds/publichost/protocol"
	"io"
	"net"
	"net/http"
	"strconv"
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
	var buffer bytes.Buffer
	request.Write(&buffer)

	r := protocol.NewRequest()
	r.Header["ContentLength"] = strconv.Itoa(buffer.Len())
	r.Body = &buffer

	t.send(r)
	hijacker, ok := response.(http.Hijacker)
	if !ok {
		panic("Response does not support hijacking")
	}
	connection, readWriter, err := hijacker.Hijack()
	if err != nil {
		panic("Hijacking failed: " + err.Error())
	}
	defer connection.Close()

	r, err = t.receive()
	if err != nil {
		panic("Could not receive request from tunnel: " + err.Error())
	}

	io.Copy(readWriter, r.Body)
}

func (t *Tunnel) send(request *protocol.Request) error {
	var buffer bytes.Buffer
	for key, value := range request.Header {
		buffer.WriteString(fmt.Sprintf("%v:%v", key, value))
	}

	_, err := buffer.WriteTo(t.connection)
	return err
}

func (t *Tunnel) receive() (*protocol.Request, error) {
	message, err := protocol.ReadRequest(t.connection)
	return message, err
}

func (t *Tunnel) Serve() {
	buffer := make([]byte, 4096)
	defer t.connection.Close()

	for {
		n, err := t.connection.Read(buffer)
		if err != nil {
			fmt.Printf("Error reading from %v: %v", t.connection.RemoteAddr(), err)
			break
		}

		msg := buffer[0:n]
		fmt.Printf("Recieved data from %v: %v\n", t.connection.RemoteAddr(), string(msg))
	}

	fmt.Printf("Closing tunnel from %v")
}

type TunnelServiceHost struct {
	tunnel *Tunnel
}

func main() {
	host := &TunnelServiceHost{}
	if err := host.ListenAndServe(); err != nil {
		fmt.Sprintf("[FATAL] Error serving PH: %v", err)
	}
}

func (t *TunnelServiceHost) ListenAndServe() error {
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

		// TODO: We need to hold tunnels, not a single tunnel.
		t.tunnel = tunnel
	}
}

func (t *TunnelServiceHost) ListenAndServeHTTP() error {
	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		if t.tunnel != nil {
			t.tunnel.HandleHttp(response, request)
		} else {
			response.WriteHeader(http.StatusNotFound)
		}
	})
	return http.ListenAndServe("", nil)
}
