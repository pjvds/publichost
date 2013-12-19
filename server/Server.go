package server

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"sync"
)

var (
	ErrAlreadyExists = errors.New("tunnel already exists")
)

type TunnelFrontEnd struct {
	Hostname string
	lock     sync.Mutex

	connection net.Conn
	readWriter *bufio.ReadWriter
}

func (t *TunnelFrontEnd) Ack() {
	r := http.Response{}
	r.StatusCode = http.StatusOK
	r.Header.Add("X-Tunnel-Hostname", t.Hostname)

	// TODO: Close tunnel on error
	writer := t.readWriter.Writer
	r.Write(writer)
	writer.Flush()
}

func (t *TunnelFrontEnd) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	t.lock.Lock()
	defer t.lock.Unlock()

	log.Debug("Proxying request %s to tunnel", request.RequestURI)

	if err := request.Write(t.readWriter); err != nil {
		log.Error("Unable to write request to tunnel: %v", err)
	}

	conn, rw, err := response.(http.Hijacker).Hijack()
	if err != nil {
		log.Error("Error hijacking tunnel response: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
	}
	defer conn.Close()

	r, err := http.ReadResponse(rw.Reader, request)
	if err != nil {
		log.Error("Unable to read response from tunnel: %v", err)
	}

	if err := r.Write(rw.Writer); err != nil {
		log.Error("Unable to write response from front end: %v", err)
	}
	if err := rw.Writer.Flush(); err != nil {
		log.Error("Unable to write response from front end: %v", err)
	}
}

type Server interface {
	ListenAndServe() error
}

func NewServer(address string) Server {
	return &server{
		tunnels: make(map[string]*TunnelFrontEnd),
	}
}

type server struct {
	address string // The TCP address to listen on
	lock    sync.Mutex
	tunnels map[string]*TunnelFrontEnd
}

func (s *server) ListenAndServe() error {
	return http.ListenAndServe(s.address, s)
}

func (s *server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	log.Debug("Incoming request: %v %v", request.Method, request.RequestURI)

	if request.URL.Host == "tunneler.publichost.me" {
		log.Debug("Handling request as new tunnel request")

		s.handleNewTunnelRequest(response, request)
		return
	} else {
		log.Debug("Handling request http request to proxy")
		s.handleHttpTraffic(response, request)
	}
}

func (s *server) handleNewTunnelRequest(response http.ResponseWriter, request *http.Request) {
	if request.Method != "CONNECT" {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	var hostname string
	if hostname = request.Header.Get("X-PublicHost-HostName"); hostname == "" {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	connection, readWriter, err := response.(http.Hijacker).Hijack()
	if err != nil {
		log.Error("Unable to hijack response: %v", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	tunnel, err := s.addNewTunnelByHostname(hostname, connection, readWriter)
	if err == ErrAlreadyExists {
		response.WriteHeader(http.StatusNotModified)
	}

	tunnel.Ack()
}

func (s *server) handleHttpTraffic(response http.ResponseWriter, request *http.Request) {
	tunnel := s.getTunnelByHostname(request.URL.Host)
	if tunnel == nil {
		response.WriteHeader(http.StatusBadGateway)
		return
	}

	tunnel.ServeHTTP(response, request)
}

func (s *server) addNewTunnelByHostname(hostname string, connection net.Conn, readWriter *bufio.ReadWriter) (*TunnelFrontEnd, error) {
	tunnel := &TunnelFrontEnd{
		Hostname:   hostname,
		connection: connection,
		readWriter: readWriter,
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.tunnels[hostname]; ok {
		return nil, ErrAlreadyExists
	}

	s.tunnels[hostname] = tunnel
	return tunnel, nil
}

// Returns the matching tunnel, or nil if not found.
func (s *server) getTunnelByHostname(hostname string) *TunnelFrontEnd {
	s.lock.Lock()
	defer s.lock.Unlock()

	if tunnel, ok := s.tunnels[hostname]; ok {
		return tunnel
	}

	return nil
}
