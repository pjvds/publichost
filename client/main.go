package main

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/handlers"
	"github.com/hashicorp/yamux"
	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{}
	app.Name = "publichost"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "publichost",
			Value:   "api.publichost.me:80",
			Usage:   "the address of the publichost server",
			EnvVars: []string{"PUBLICHOST"},
		},
		&cli.BoolFlag{
			Name: "tls",
		},
	}
	app.Commands = []*cli.Command{
		&cli.Command{
			Name: "dir",
			Action: func(ctx *cli.Context) error {
				localDir := ctx.Args().First()
				if len(localDir) == 0 {
					log.Fatal("local directory not specified")
				}

				if _, err := os.Stat(localDir); err != nil {
					log.Fatal(err.Error())
				}

				log.Println("connecting to server")

				var conn net.Conn
				var err error
				if ctx.Bool("tls") {
					conn, err = tls.Dial("tcp", ctx.String("publichost"), nil)
				} else {
					conn, err = net.Dial("tcp", ctx.String("publichost"))
				}

				if err != nil {
					log.Fatal(err.Error())
				}
				if _, err = conn.Write([]byte("GET /tunnel HTTP/1.1\r\nHost: api.publichost.io\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==\r\nSec-WebSocket-Version: 13\r\n\r\n")); err != nil {
					log.Fatal(err.Error())
				}

				log.Println("opening tunnel")
				reader := bufio.NewReader(conn)
				response, err := http.ReadResponse(reader, nil)
				if err != nil {
					log.Fatal(err.Error())
				}

				log.Print("tunnel available at: " + response.Header.Get("X-Publichost-Address"))
				tunnel, err := yamux.Server(conn, nil)
				if err != nil {
					log.Fatal(err.Error())
				}

				fileserver := http.FileServer(http.Dir(localDir))
				handler := handlers.CombinedLoggingHandler(os.Stdout, fileserver)
				if err := http.Serve(tunnel, handler); err != nil {
					log.Fatal(err.Error())
				}

				return nil
			},
		},
		&cli.Command{
			Name: "http",
			Action: func(ctx *cli.Context) error {
				localUrl := ctx.Args().First()

				log.Println("connecting to server")
				conn, err := tls.Dial("tcp", "api.publichost.io:443", nil)
				if err != nil {
					log.Fatal(err.Error())
				}

				request, err := http.NewRequest("GET", "api.publichost.io", nil)
				if err != nil {
					log.Fatal(err)
				}

				request.Header.Set("X-Publichost-Local", localUrl)
				request.Header.Set("Upgrade", "websocket")
				request.Header.Set("Connection", "Upgrade")
				request.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")
				request.Header.Set("Sec-WebSocket-Version", "13")

				if err = request.Write(conn); err != nil {
					log.Fatal(err.Error())
				}

				log.Println("opening tunnel")
				reader := bufio.NewReader(conn)
				response, err := http.ReadResponse(reader, nil)
				if err != nil {
					log.Fatal(err.Error())
				}

				log.Print("tunnel available at: " + response.Header.Get("X-Publichost-Address"))
				tunnel, err := yamux.Server(conn, nil)
				if err != nil {
					log.Fatal(err.Error())
				}

				local, err := url.Parse(localUrl)
				if err != nil {
					log.Fatal(err.Error())
				}

				handler := handlers.CombinedLoggingHandler(os.Stdout, httputil.NewSingleHostReverseProxy(local))

				if err := http.Serve(tunnel, handler); err != nil {
					log.Fatal(err.Error())
				}

				return nil
			},
		},
	}
	app.Run(os.Args)
}
