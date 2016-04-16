package main

import(
    "github.com/codegangsta/cli"
    "net"
    "github.com/hashicorp/yamux"
    "log"
    "bufio"
    "net/http"
    "net/http/httputil"
    "net/url"
)

func main() {
    app := cli.NewApp()
    app.Name = "publichost"
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name: "url",
            Value: "http://localhost:3000",
            Usage: "the local url to make publicly available",
            EnvVar: "URL",
          },
          cli.StringFlag{
              Name: "publichost, p",
              Value: "localhost:5000",
              Usage: "the address of the publichost server",
              EnvVar: "PUBLICHOST",
          },
    }
    app.Action = func(ctx *cli.Context) {
        localUrl := ctx.String("url")
        server := ctx.String("publichost")

        log.Println("connecting to server")
        conn, err := net.Dial("tcp", server)
        if err != nil {
            log.Fatal(err.Error())
        }

        log.Println("opening tunnel")
        reader := bufio.NewReader(conn)
        line, err := reader.ReadString('\n')
        if err != nil {
            log.Fatal(err.Error())
        }

        log.Print("tunnel available at: " + line)
        tunnel, err := yamux.Server(conn, nil)
        if err != nil {
            log.Fatal(err.Error())
        }

        local, err := url.Parse(localUrl)
        if err != nil {
            log.Fatal(err.Error())
        }

        if err := http.Serve(tunnel, httputil.NewSingleHostReverseProxy(local)); err != nil {
            log.Fatal(err.Error())
        }
    }
    app.RunAndExitOnError()
}
