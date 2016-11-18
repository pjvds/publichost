package main

import (
	"bufio"
	"crypto/tls"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gorilla/handlers"
	"github.com/hashicorp/yamux"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "publichost"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "publichost, p",
			Value:  "",
			Usage:  "the address of the publichost server",
			EnvVar: "PUBLICHOST",
		},
	}
	app.Commands = []cli.Command{
		cli.Command{
			Name: "dir",
			Action: func(ctx *cli.Context) {
				localDir := ctx.Args().First()

				log.Println("connecting to server")
				conn, err := tls.Dial("tcp", "api.publichost.io:443", nil)
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

				handler := handlers.CombinedLoggingHandler(os.Stdout, http.FileServer(http.Dir(localDir)))
				if err := http.Serve(tunnel, handler); err != nil {
					log.Fatal(err.Error())
				}
			},
		},
		cli.Command{
			Name: "http",
			Action: func(ctx *cli.Context) {
				localUrl := ctx.Args().First()

				log.Println("connecting to server")
				conn, err := tls.Dial("tcp", "api.publichost.io:443", nil)
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

				local, err := url.Parse(localUrl)
				if err != nil {
					log.Fatal(err.Error())
				}

				handler := handlers.CombinedLoggingHandler(os.Stdout, httputil.NewSingleHostReverseProxy(local))

				if err := http.Serve(tunnel, handler); err != nil {
					log.Fatal(err.Error())
				}
			},
		},
	}
	app.RunAndExitOnError()
}
