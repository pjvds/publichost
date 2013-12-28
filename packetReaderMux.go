package publichost

import (
	"fmt"
	"github.com/pjvds/publichost/network"
)

type packetReaderMux struct {
	conn network.Connection

	Requests  chan *Request
	Responses chan *Response
}

func newPacketReaderMux(conn network.Connection) *packetReaderMux {
	return &packetReaderMux{
		conn:      conn,
		Requests:  make(chan *Request, 10),
		Responses: make(chan *Response, 10),
	}
}

func (p *packetReaderMux) Read() (err error) {
	var packet *network.Packet
	var request *Request
	var response *Response

	for {
		if packet, err = p.conn.Receive(); err != nil {
			return
		}

		typeId := packet.TypeId()
		switch {
		case typeId == TRequest:
			if request, err = ReadRequest(packet.CreateContentReader()); err != nil {
				return
			}
			p.Requests <- request
		case typeId == TResponse:
			if response, err = ReadResponse(packet.CreateContentReader()); err != nil {
				return
			}
			p.Responses <- response
		default:
			return fmt.Errorf("unknown packet id %#h", p)
		}
	}
}
