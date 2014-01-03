package server

import (
	"bufio"
	"github.com/pjvds/publichost/tunnel"
	"net"
	"net/http"
)

type PublicHostServer interface {
	Serve() error
}

type publicHostServer struct {
	listener net.Listener

	tunnels map[string]tunnel.Tunnel
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
	err := make(chan error)

	go func() {
		if e := p.serveTunnels(); e != nil {
			err <- e
		}
	}()

	go func() {
		if e := p.serveHttp(); e != nil {
			err <- e
		}
	}()

	return <-err
}

func (p *publicHostServer) serveTunnels() error {
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

func (p *publicHostServer) serveHttp() error {
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		log.Debug("incomming http request from %v to %v", req.RemoteAddr, req.RequestURI)

		if t, ok := p.tunnels[req.URL.Host]; ok {
			// TODO: We need to rewrite the destination
			id, err := t.OpenStream("tcp", "127.0.0.1:4000")
			if err != nil {
				log.Error("error opening stream: %v", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			stream := tunnel.NewTunneledStream(id, t)
			req.Write(stream)

			bufReader := bufio.NewReader(stream)

			response, err := http.ReadResponse(bufReader, req)
			if err != nil {
				log.Error("error reading response from tunneled stream: %v", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			conn, readWriter, _ := rw.(http.Hijacker).Hijack()
			defer conn.Close()

			if err := response.Write(readWriter); err != nil {
				log.Error("error writing response to remote: %v", err)
				return
			}

			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	})

	return http.ListenAndServe("0.0.0.0:4000", nil)
}
