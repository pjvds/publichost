package main

import(
    "strconv"
    "strings"
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

type TunnelSession struct {
    id int
    session *yamux.Session
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

    for id := 0; ; id++{
         conn, err := listener.Accept()
         if err != nil {
             return err
         }
         log.Printf("connection accepted %v\n", conn.RemoteAddr())

         // perform handshake
         go func(conn net.Conn, id int) {
            hostname := fmt.Sprintf("%v.%v", id, publicHostname)
            publicAddress := fmt.Sprintf("http://%v", hostname)

            if _, err := conn.Write([]byte(publicAddress + "\n")); err != nil {
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


        accepted := make(chan TunnelSession)
        go Accept(accepted, "localhost:8080", listener)

        var tunnelsLock sync.RWMutex
        tunnels := make(map[int]TunnelSession)

        go http.ListenAndServe(httpAddress, http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
            if !strings.HasSuffix(request.Host, ".localhost:8080") {
                log.Println("invalid hostname")
                response.Write([]byte("<html><body>invalid hostname</body></html>"))
                return
            }

            id, err := strconv.Atoi(strings.TrimSuffix(request.Host, ".localhost:8080"))
            if err != nil {
                log.Println("invalid tunnel id: " + err.Error())
                response.Write([]byte("<html><body>invalid tunnel id: "+err.Error()+"</body></html>"))
                return
            }

            tunnelsLock.RLock()
            tunnel, ok := tunnels[id]
            tunnelsLock.RUnlock()

            if !ok {
                log.Printf("no tunnel with id %v", id)
                response.Write([]byte("<html><body>no session found with id <strong>" + strconv.Itoa(id) + "</strong></body></html>"))
                return
            }

            tunnel.ServeHTTP(response, request)
        }))

        for session := range accepted {
            tunnelsLock.Lock()
            tunnels[session.id] = session
            tunnelsLock.Unlock()
        }
    }
    app.RunAndExitOnError()
}
