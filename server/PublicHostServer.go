package server

import (
	"github.com/pjvds/publichost/tunnel"
	"github.com/pjvds/publichost/net/message"
	pnet "github.com/pjvds/publichost/net"
	"net"
	"net/http"
	"bufio"
	"math/rand"
	"fmt"
	"io"
)

type publicHostServer struct {
	listener net.Listener
	hostname string

	tunnels map[string]tunnel.Tunnel
	mux *http.ServeMux
}

func NewServer(address, hostname string) (server http.Handler, err error) {
	s := publicHostServer{
		hostname: hostname,
		tunnels: make(map[string]tunnel.Tunnel),
		mux: http.NewServeMux(),
	}
	s.mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		// All requests marked with X-PUBLICHOST header
		if req.Header.Get("X-PUBLICHOST") == "true" ||
		   req.Header.Get("X-PUBLICHOST") == "1" {
			s.handleConnectRequests(rw, req)
			return
		}

		s.handlePotentialTunnelRequest(rw, req)
		return
	})

	server = s.mux
	return 
}

func ListenAndServe(address string, hostname string) (err error) {
	var server http.Handler
	if server, err = NewServer(address, hostname); err != nil {
		return
	}
	return http.ListenAndServe(address, server)
}

func (p *publicHostServer) handleConnectRequests(rw http.ResponseWriter, req *http.Request) {
	log.Debug("handling connect request")
	if req.Header.Get("X-PUBLICHOST") == "" {
		http.Error(rw, "Missing header: X-PUBLICHOST", http.StatusBadRequest)
		return
	}
	if req.Method != "CONNECT" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
	}

	conn, bufRW, err := rw.(http.Hijacker).Hijack()
	if err != nil {
		log.Error("cannot hijack connection: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	response := &http.Response{
		Status: "200 OK",
		StatusCode: 200,
		Proto: "HTTP/1.0",
	}

	if err = response.Write(bufRW); err != nil {
		log.Debug("cannot write response to buffer: %v", err)
		return
	}
	if err = bufRW.Flush(); err != nil {
		log.Debug("cannot flush response: %v", err)
		return
	}
	log.Debug("confirmed request")

	p.serveTunnel(conn, bufRW)
}

func (p *publicHostServer) serveTunnel(conn net.Conn, bufRW *bufio.ReadWriter) {
	reader := message.NewReader(bufRW.Reader)
	writer := message.NewWriter(bufRW.Writer)

	log.Debug("waiting for open tunnel request")

	m, err := reader.Read()
	if err != nil {
		log.Debug("error reading: %v", err)
	}
	if m.TypeId != message.OpOpenTunnel {
		log.Debug("first request was not open tunnel, but: %v", m.String())
		conn.Close()
		return
	}

	log.Debug("open tunnel request received")
	
	name := rand_str(3)
	hostname := fmt.Sprintf("%v.%v", name, p.hostname)

	if err := writer.Write(message.NewMessage(message.Ack, m.CorrelationId, []byte(hostname))); err != nil {
		conn.Close()
		return
	}

	p.tunnels[hostname] = tunnel.NewFrondend(pnet.NewClientConnection(conn))
	log.Info("opened new tunnel at: %v", hostname)
}

func rand_str(length int) string {
    alphanum := "abcdefghijklmnopqrstuvwxyz"
    result := make([]byte, length, length)

    for i := 0; i < length; i++ {
    	r := rand.Intn(len(alphanum))
    	result[i] = alphanum[r]
    }
    return string(result)
}

func (p *publicHostServer) handlePotentialTunnelRequest(rw http.ResponseWriter, req *http.Request) {
	var ok bool
	var t tunnel.Tunnel

	if t, ok = p.tunnels[req.Host]; !ok {
		log.Debug("no tunnel found for host: %v", req.Host)
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	// TODO: We need to rewrite the destination
	id, err := t.OpenStream("tcp", "127.0.0.1:4000")
	if err != nil {
		log.Error("error opening stream: %v", err)
		http.Error(rw, err.Error(), http.StatusRequestTimeout)
		return
	}
	defer t.CloseStream(id)
	log.Debug("opened stream\n")

	stream := tunnel.NewTunneledStream(id, t)
	if err = req.Write(stream); err != nil {
		log.Error("cannot write request: %v", err)
		http.Error(rw, err.Error(), http.StatusRequestTimeout)
		return
	}

	response, err := http.ReadResponse(bufio.NewReader(stream), req)
	if err != nil {
		log.Error("cannot read response: %v", err)
		http.Error(rw, err.Error(), http.StatusRequestTimeout)
		return
	}

	for k, _ := range response.Header {
		rw.Header().Set(k, response.Header.Get(k))
	}

	log.Debug("read response")
	if _, err := io.Copy(rw, response.Body); err != nil {
		log.Error("cannot copy body: %v", err)
		return
	}

	log.Debug("done")
}