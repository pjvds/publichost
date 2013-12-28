package publichost

import (
	"github.com/pjvds/publichost/network"
	"time"
)

type RequestHandler func(request *Request)
type ResponseHandler func(response *Response)

const (
	TRequest  = byte(iota)
	TResponse = byte(iota)
)

const (
	OpOpenTunnel  = byte(iota)
	OpCloseTunnel = byte(iota)
	OpKeepAlive   = byte(iota)
)

const (
	StatusOK    = byte(iota)
	StatusError = byte(iota)
)

type TunnelFrontEnd struct {
	conn network.Connection

	mux *packetReaderMux
}

func newTunnelFrondEnd(conn network.Connection) *TunnelFrontEnd {
	return &TunnelFrontEnd{
		conn: conn,
		mux:  newPacketReaderMux(conn),
	}
}

func (t *TunnelFrontEnd) Serve(s *Server) {
	defer t.close()

	select {
	case request := <-t.mux.Requests:
		log.Info("request received: %s", request)
	case response := <-t.mux.Responses:
		log.Info("response received: %s", response)
	}
}

func (t *TunnelFrontEnd) close() {
	if err := t.conn.Close(); err != nil {
		log.Warning("error closing tunnel connection: %v", err)
	}

	log.Info("closed tunnel %v", t.conn)
}

func (t *TunnelFrontEnd) handshake() error {
	// TODO: Implement protocol negociation
	return nil
}

func (t *TunnelFrontEnd) Close() error {
	return t.conn.Close()
}

type Server struct {
}

func (s *Server) Serve() {

}

func (s *Server) serveTunnelConnections(listener network.Listener) (err error) {
	for {
		var conn network.Connection
		if conn, err = listener.Accept(); err != nil {
			if network.IsTemporaryError(err) {
				duration := 1 * time.Second
				log.Info("temporary accept error: %v; retry in %v", err, duration)

				time.Sleep(duration)
				continue
			}
			return err
		}

		log.Info("accepted connection from %v", conn)

		tunnel := newTunnelFrondEnd(conn)
		go tunnel.Serve(s)
	}
}
