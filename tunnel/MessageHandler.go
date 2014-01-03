package tunnel

import (
	"github.com/pjvds/publichost/net/message"
)

type MessageHandler interface {
	Handle(response ResponseWriter, request *message.Message) error
}
