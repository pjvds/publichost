package main

import (
	"flag"
	"github.com/pjvds/publichost/server"
    "github.com/yvasiyarov/gorelic"
)

var (
	address = flag.String("address", "0.0.0.0:8080", "The address to bind to")
    hostname = flag.String("hostname", "publichost.me", "The hostname of the frontend")
    newrelic = flag.String("newrelic", "", "The new relic license key")
)

func main() {
	flag.Parse()

    if *newrelic != "" {
        agent := gorelic.NewAgent()
        agent.Verbose = false
        agent.NewrelicLicense = *newrelic
        agent.Run()
    }

	println("publichost - v0.1")
	println("hosting at: " + *address)

	if err := server.ListenAndServe(*address, *hostname); err != nil {
		println("err: " + err.Error())
	}
}
