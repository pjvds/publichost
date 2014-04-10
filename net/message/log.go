package message

import (
    "github.com/op/go-logging"
)

var (
    log *logging.Logger
)

func init() {
    log = logging.MustGetLogger("message")
}
