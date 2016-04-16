package main

import(
    "net"
    "net/http"
    "github.com/codegangsta/cli"
    "github.com/hashicorp/yamux"
    "log"
    "fmt"
    "bufio"
    "io"
    "sync"
)

func main() {
    app := cli.NewApp()
    app.Name = "publichost server"
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name: "api",
            Value: ":5000",
            Usage: "the api address to bind to serve api",
            EnvVar: "API",
          },
        cli.StringFlag{
            Name: "http",
            Value: ":8080",
            Usage: "the address to bind to serve http",
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

        tunnels := make(map[string]*yamux.Session)
        var tunnelsLock sync.RWMutex

        go http.ListenAndServe(httpAddress, http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
            tunnelsLock.RLock()

            id := request.Host
            session, ok := tunnels[id]
            tunnelsLock.RUnlock()
            
            if !ok {
                response.Write([]byte("<html><body>no session found with id <strong>" + id + "</strong></body></html>"))
                tunnelsLock.RUnlock()
                return
            }

            conn, err := session.Open()
            if err != nil {
                log.Println(err.Error())
                response.Write([]byte("<html><body>" + err.Error() + "</body></html>"))
                return
            }

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
        }))

        id := 0
        for {
            id++
            conn, err := listener.Accept()
            if err != nil {
                log.Fatal(err.Error())
            }

            go func(conn net.Conn, id int) {
                hostname := fmt.Sprintf("%v.%v", id, httpAddress)
                remoteAddress := fmt.Sprintf("http://%v", hostname)

                if _, err := conn.Write([]byte(remoteAddress + "\n")); err != nil {
                    log.Println(err.Error())
                    return
                }

                    session, err :=yamux.Client(conn, nil)
                    if err != nil {
                        log.Println(err.Error())
                        return
                    }

                    tunnelsLock.Lock()
                    tunnels[hostname] = session
                    tunnelsLock.Unlock()

                    log.Println(hostname + " tunnel active")
            }(conn, id)
        }
    }
    app.RunAndExitOnError()
}
