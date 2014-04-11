package main

import (
	"github.com/op/go-logging"
)

var (
	log *logging.Logger
)

func init() {
    //logging.SetLevel(logging.INFO, "")
	log = logging.MustGetLogger("main")
}

