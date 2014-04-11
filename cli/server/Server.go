package main

import (
	"flag"
	"github.com/pjvds/publichost/server"
    "github.com/yvasiyarov/gorelic"
    "net/http"
)

var (
	address = flag.String("address", "0.0.0.0:8080", "The address to bind to")
    hostname = flag.String("hostname", "publichost.me", "The hostname of the frontend")
    newrelic = flag.String("newrelic", "", "The new relic license key")
)

func main() {
	flag.Parse()

    var agent *gorelic.Agent
    if *newrelic != "" {
        agent = gorelic.NewAgent()
        agent.Verbose = true
        agent.NewrelicLicense = *newrelic
        agent.NewrelicName = "publichost server"
        agent.Run()
    }

	println("publichost - v0.1")
	println("hosting at: " + *address)

    s, err := server.NewServer(*address, *hostname)
	if err != nil {
		println("err: " + err.Error())
	}

    if agent != nil {
        s = agent.WrapHTTPHandler(s)
    }

    http.ListenAndServe(*address, s)
}
