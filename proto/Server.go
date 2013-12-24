package proto

import (
	"net"
)

type TunnelFrontEndConnection struct {
	// Signal that this tunnel is closed.
	closed chan *TunnelFrontEndConnection
}


type Server struct {
	tunnels map[string]*TunnelFrontEndConnection

	// Signals when a tunnel got closed.
	closed chan *TunnelFrontEndConnection
}

func (s *Server) serve() {
	for {
            switch{
        case t := s.closed

        }
	}
}
