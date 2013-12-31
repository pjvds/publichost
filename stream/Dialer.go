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
	return net.Dial(network, address)
}
