package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/dustinkirkland/golang-petname"
	"github.com/hashicorp/yamux"
)

type TunnelSession struct {
	id            int
	name          string
	session       *yamux.Session
	remoteAddress string
}

func (this TunnelSession) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	conn, err := this.session.Open()
	if err != nil {
		log.Println(err.Error())
		response.Write([]byte("<html><body>" + err.Error() + "</body></html>"))
		return
	}
	defer conn.Close()

	if err := request.Write(conn); err != nil {
		log.Println(err.Error())
		return
	}

	tunnelResponse, err := http.ReadResponse(bufio.NewReader(conn), request)
	if err != nil {
		response.Write([]byte("<html><body>" + err.Error() + "</body></html>"))
		return
	}

	for header, values := range tunnelResponse.Header {
		for _, value := range values {
			response.Header().Add(header, value)
		}
	}

	response.WriteHeader(tunnelResponse.StatusCode)
	if tunnelResponse.Body != nil {
		io.Copy(response, tunnelResponse.Body)
	}
}

// Accepts new tunnel session in a blocking fashion. When it returns it closes the listener
// and returns with the error why it returned.
func Accept(accepted chan TunnelSession, publicHostname string, listener net.Listener) error {
	defer listener.Close()

	for id := 0; ; id++ {

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("FATAL: %v", err)
			return err
		}

		log.Printf("connection accepted %v\n", conn.RemoteAddr())

		// perform handshake
		go func(conn net.Conn, id int) {
			name := petname.Generate(2, "-")

			hostname := fmt.Sprintf("%v.%v", name, publicHostname)
			publicAddress := fmt.Sprintf("http://%v", hostname)

			reader := bufio.NewReader(conn)
			_, err := http.ReadRequest(reader)
			if err != nil {
				log.Println(err.Error())
				return
			}

			if _, err := conn.Write([]byte("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=\r\nX-Publichost-Address: " + publicAddress + "\r\n\r\n")); err != nil {
				log.Println(err.Error())
				return
			}

			session, err := yamux.Client(conn, nil)
			if err != nil {
				log.Println(err.Error())
				return
			}

			log.Printf("tunnel created %v->%v\n", publicAddress, conn.RemoteAddr())
			accepted <- TunnelSession{
				id,
				name,
				session,
				publicAddress,
			}
		}(conn, id)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "publichost server"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api",
			Value:  ":5000",
			Usage:  "the api address to bind to serve api",
			EnvVar: "API",
		},
		cli.StringFlag{
			Name:   "http",
			Value:  ":8080",
			Usage:  "the address to bind to serve http",
			EnvVar: "HTTP",
		},
	}
	app.Action = func(ctx *cli.Context) {
		apiAddress := ctx.String("api")
		httpAddress := ctx.String("http")
		listener, err := net.Listen("tcp", apiAddress)
		if err != nil {
			log.Fatal(err.Error())
		}

		accepted := make(chan TunnelSession)
		go Accept(accepted, "publichost.me", listener)

		var tunnelsLock sync.RWMutex
		tunnels := make(map[string]TunnelSession)

		subdomain := regexp.MustCompile("[A-Za-z0-9](?:[A-Za-z0-9\\-]{0,61}[A-Za-z0-9])?")

		go http.ListenAndServe(httpAddress, http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			name := subdomain.FindString(request.Host)
			log.Printf("handling incoming request %v->%v\n", request.Host, name)
			if len(name) == 0 {
				log.Println("missing tunnel name")
				response.Write([]byte("<html><body>missing tunnel name</body></html>"))
				return
			}

			tunnelsLock.RLock()
			tunnel, ok := tunnels[name]
			tunnelsLock.RUnlock()

			if !ok {
				log.Printf("no tunnel with name: %v", name)
				response.Write([]byte("<html><body>no session found with name <strong>" + name + "</strong></body></html>"))
				return
			}

			tunnel.ServeHTTP(response, request)
		}))

		for session := range accepted {
			tunnelsLock.Lock()
			tunnels[session.name] = session
			tunnelsLock.Unlock()
		}
	}
	app.RunAndExitOnError()
}
