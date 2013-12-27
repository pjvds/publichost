package publichost

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/pjvds/publichost/network"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

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

type TunnelConnection struct {
	conn net.Conn

	server *Server
}

func (t *TunnelConnection) Serve() {
	defer t.close()

	reader := bufio.NewReader(t.conn)
	writer := bufio.NewWriter(t.conn)

	if err := t.handshake(reader, writer); err != nil {
		log.Info("unable to handshake: %v; tunnel will close", err)
		return
	}

	for {
		var request *Request
		var err error

		if request, err = readNextRequest(reader); err != nil {
			log.Notice("error reading request: %v; tunnel will close", err)
			break
		}

		if request.Type == OpCloseTunnel {
			log.Notice("close request received")
		}
	}
}

func (t *TunnelConnection) handshake(reader *bufio.Reader, writer *bufio.Writer) (err error) {
	var request *Request
	if request, err = readNextRequest(reader); err != nil {
		return
	}

	if request.Type != OpOpenTunnel {
		response := NewResponse(request.Id, StatusError, bytes.NewBufferString("invalid request type"))
		response.Write(writer)
		writer.Flush()

		log.Info("invalid request type received for handshake: %v; expected: %v", request.Type, OpOpenTunnel)
		return errors.New("invalid request type received")
	}

	response := NewResponse(request.Id, StatusOK, nil)
	if err = response.Write(writer); err != nil {
		return
	}
	if err = writer.Flush(); err != nil {
		return
	}
	return
}

// Read the next request. When an temporary error occurs
// it will retry until a request is received. When a non
// temporary error occurs, this will returned.
func readNextRequest(reader io.Reader) (request *Request, err error) {
	for {
		if request, err = ReadRequest(reader); err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				toSleep := 1 * time.Second

				log.Warning("temporary error reading request: %v; retry in ", toSleep)
				time.Sleep(toSleep)
				continue
			}
			return
		}
	}
}

// Read the next response. When an temporary error occurs
// it will retry until a response is received. When a non
// temporary error occurs, this will returned.
func readNextResponse(reader io.Reader) (response *Response, err error) {
	for {
		if response, err = ReadResponse(reader); err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				toSleep := 1 * time.Second

				log.Warning("temporary error reading response: %v; retry in ", toSleep)
				time.Sleep(toSleep)
				continue
			}
			return
		}
	}
}

func (t *TunnelConnection) close() {
	defer t.conn.Close()

	log.Info("closing tunnel with client %v", t.conn.RemoteAddr())
}

type Server struct {
	serveWaitGroup sync.WaitGroup

	commandAddress string
	dataAddress    string
}

func ListenAndServe(commandAddress, dataAddress string) error {
	s := &Server{
		commandAddress: commandAddress,
		dataAddress:    dataAddress,
	}
	return s.Serve()
}

func (s *Server) Serve() (err error) {
	var commandListener
	var dataListener net.Listener
	if commandListener, err = net.Listen("tcp", s.commandAddress); err != nil {
		log.Error("error opening command listener at %v: %v", s.commandAddress, err)
		return
	}
	defer commandListener.Close()

	if dataListener, err = net.Listen("tcp", s.dataAddress); err != nil {
		log.Error("error opening data listener at %v: %v", s.commandAddress, err)
		return
	}
	defer dataListener.Close()

	return s.serve(commandListener, dataListener)
}

func (s *Server) serve(commandListener, dataListener net.Listener) error {
	errored := make(chan error)
	go func() {
		if err := s.serveCommands(commandListener); err != nil {
			log.Error("command listener stopped: %v", err)
			errored <- err
			return
		}

		log.Info("Command listener finished")
	}()

	go func() {
		if err := s.serveData(dataListener); err != nil {
			log.Error("data listener stopped: %v", err)
			errored <- err
			return
		}

		log.Info("Command listener finished")
	}()

	err := <-errored

	return err
}

func (s *Server) serveCommands(listener net.Listener) error {
	defer s.serveWaitGroup.Done()
	defer listener.Close()

	s.serveWaitGroup.Add(1)

	log.Debug("waiting for incomming command connections")

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
		// a new goroutine and continue accepting connections.
		log.Info("accepted connection from %v", conn.RemoteAddr())
		tunnel := &TunnelConnection{
			conn:   conn,
			server: s,
		}
		go tunnel.Serve()
	}
}

func (s *Server) serveData(listener net.Listener) error {
	defer s.serveWaitGroup.Done()
	s.serveWaitGroup.Add(1)

	log.Debug("waiting for incomming data connections")

	mux := http.NewServeMux()
	return http.Serve(listener, mux)
}
