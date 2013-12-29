package server

import (
	"net"
)

type StreamDialer interface {
	Dial(network, address string) (conn StreamConnection, err error)
}

type networkStreamDialer struct {
}

func newNetworkStreamDialer() *networkStreamDialer {
	return &networkStreamDialer{}
}

func (d *networkStreamDialer) Dial(network, address string) (conn StreamConnection, err error) {
	return net.Dial(network, address)
}
