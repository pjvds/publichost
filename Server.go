package publichost

import (
	"net"
	"time"
)

type Response struct {
	RequestId uint16
}

type TunnelConnection struct {
	conn net.Conn

	server *Server
}

func (t *TunnelConnection) Serve() {
	defer t.conn.Close()

	for {

	}
}

type Server struct {
}

func (s *Server) serve(listener net.Listener) error {
	defer listener.Close()

	for {
		// Accept a new connection.
		conn, err := listener.Accept()

		// If there is an error, see if it is tempoary and retry
		// or fail by return the error.
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				toSleep := 1 * time.Second
				log.Warning("temporary accept error: %v; retry in %v", ne.Error(), toSleep.String())

				time.Sleep(toSleep)
				continue
			}

			log.Error("accept error: %v", err.Error())
			return err
		}

		// We successfully accepted an connection. Serve it in
		// a new goroutine and accept a accept a new connection.
		log.Info("accepted connection from %v", conn.RemoteAddr())
		tunnel := &TunnelConnection{
			conn:   conn,
			server: s,
		}
		go tunnel.Serve()
	}
}
