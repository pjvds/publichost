package publichost

import (
	"net"
)

type StreamId uint32

type OutgoingData struct {
    Data []byte

    Ack chan bool
    Err chan error
}

type IncommingData struct {
    Data []byte

    Ack chan bool
    Err chan error
}

type StreamFrondEnd struct {
	Id   StreamId
	conn net.Conn

    outgoing chan *OutgoingData
    incomming chan *
}

type client struct {
    streamManager 
}
