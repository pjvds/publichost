package publichost

import (
	"net"
    "bufio"
    "io"
)

type streamId uint32

type outgoingData struct {
    Data []byte

    Ack chan bool
    Err chan error
}

type Operation struct {
    Data []byte

    Ack chan bool
    Err chan error
}

type streamFrondEnd struct {
	Id   StreamId
	conn net.Conn

    outgoing chan *OutgoingData
    incomming chan *
}

type Connection interface{
    io.ReadWriteCloser
    String() string
}

type clientMux struct {
    
}

func (c *clientMux) Handle(message Message) {
    
}
