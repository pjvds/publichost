package server

import (
	"github.com/pjvds/publichost/tunnel"
	"net"
)

type PublicHostServer interface {
	Serve() error
}

type publicHostServer struct {
	listener net.Listener
}

func ListenAndServe(address string) (err error) {
	var addr *net.TCPAddr
	var listener net.Listener

	if addr, err = net.ResolveTCPAddr("tcp", address); err != nil {
		return
	}

	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		return
	}

	server := publicHostServer{
		listener: listener,
	}
	return server.Serve()
}

func (p *publicHostServer) Serve() error {
	defer p.listener.Close()

	for {
		conn, err := p.listener.Accept()
		if err != nil {
			return err
		}

		tunnelHost := tunnel.NewTunnelHost(conn)
		go tunnelHost.Serve()
	}
}
