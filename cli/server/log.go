package main

import (
    "os"
	"github.com/op/go-logging"
)

var (
	log *logging.Logger
)

func init() {
    // Setup one stdout and one syslog backend.
    logBackend := logging.NewLogBackend(os.Stderr, "", 0)
    logBackend.Color = true

    logging.SetBackend(logBackend)
	log = logging.MustGetLogger("main")
}

