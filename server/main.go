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

	petname "github.com/dustinkirkland/golang-petname"
	"github.com/rs/xid"

	"github.com/codegangsta/cli"
	"github.com/hashicorp/yamux"
)

var nameExp = regexp.MustCompile("(?P<name>.*?)\\.")

type Tunnel struct {
	id           xid.ID
	name         string
	hostname     string
	session      *yamux.Session
	localAddress string
}

func (this Tunnel) ServeHTTP(response http.ResponseWriter, request *http.Request) {
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
func Accept(accepted chan Tunnel, publicHostname string, listener net.Listener) error {
	defer listener.Close()

	var id xid.ID
	for {
		id = xid.New()

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("FATAL: %v", err)
			return err
		}

		log.Printf("connection accepted %v\n", conn.RemoteAddr())

		// perform handshake
		go func(conn net.Conn, id xid.ID) {
			name := petname.Generate(2, "-")

			hostname := fmt.Sprintf("%v.%v", name, publicHostname)
			publicAddress := fmt.Sprintf("http://%v", hostname)

			reader := bufio.NewReader(conn)
			request, err := http.ReadRequest(reader)
			if err != nil {
				log.Println(err.Error())
				return
			}

			localAddress := request.Header.Get("X-Publichost-Local")
			session, err := yamux.Client(conn, nil)
			if err != nil {
				log.Println(err.Error())
				return
			}

			response := http.Response{
				Status:     "200 OK",
				StatusCode: 200,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: http.Header{
					"X-Publichost-Address": []string{hostname},
				},
			}

			if err := response.Write(conn); err != nil {
				log.Println(err.Error())
				return
			}

			log.Printf("tunnel created %v->%v\n", publicAddress, localAddress)
			accepted <- Tunnel{
				id,
				name,
				hostname,
				session,
				localAddress,
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

		accepted := make(chan Tunnel)
		go Accept(accepted, "publichost.me", listener)

		var tunnelsLock sync.RWMutex
		tunnels := make(map[string]Tunnel)

		go http.ListenAndServe(httpAddress, http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			matches := nameExp.FindStringSubmatch(request.Host)
			if len(matches) <= 1 {
				log.Println("failed to match tunnel name from host: %v", request.Host)
				response.Write([]byte("<html><body>failed to match tunnel name</body></html>"))
				return
			}

			log.Println(matches)

			name := matches[1]

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
