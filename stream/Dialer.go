package stream

import (
	"net"
)

type Dialer interface {
	Dial(network, address string) (s Stream, err error)
}

type netDialer struct {
}

func NewDialer() Dialer {
	return &netDialer{}
}

func (n *netDialer) Dial(network, address string) (s Stream, err error) {
	var conn *net.TCPConn
    var raddr *net.TCPAddr

    if raddr, err = net.ResolveTCPAddr(network, address); err != nil {
        return
    }

    if conn, err = net.DialTCP("tcp", nil, raddr); err != nil {
        return
    }

    s = conn
    return
}
