package server

import (
	"bufio"
	"github.com/pjvds/publichost/tunnel"
	"github.com/pjvds/publichost/net/message"
	pnet "github.com/pjvds/publichost/net"
	"net"
	"net/http"
	"bytes"
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
		tunnels: make(map[string]tunnel.Tunnel),
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

		reader := message.NewReader(conn)
		writer := message.NewWriter(conn)

		m, err := reader.Read()
		if err != nil || m.TypeId != message.OpOpenTunnel {
			conn.Close()
			continue
		}

		if err := writer.Write(message.NewMessage(message.Ack, m.CorrelationId, nil)); err != nil {
			conn.Close()
			continue
		}

		p.tunnels["foobar"] = tunnel.NewFrondend(pnet.NewClientConnection(conn))
	}
}

func (p *publicHostServer) serveHttp() error {
	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		log.Debug("incomming http request from %v to %v", req.RemoteAddr, req.RequestURI)

		if t, ok := p.tunnels["foobar"]; ok {
			// TODO: We need to rewrite the destination
			id, err := t.OpenStream("tcp", "127.0.0.1:4000")
			if err != nil {
				log.Error("error opening stream: %v", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer t.CloseStream(id)
			log.Debug("opened stream\n")

			stream := tunnel.NewTunneledStream(id, t)
			var buffer bytes.Buffer
			req.Write(&buffer)

			if _, err = stream.Write(buffer.Bytes()); err != nil {
				log.Debug("error written stream: %v", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			log.Debug("written stream")

			bufReader := bufio.NewReader(stream)
			response, err := http.ReadResponse(bufReader, req)
			if err != nil {
				log.Error("error reading response from tunneled stream: %v", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Debug("read response")

			conn, readWriter, _ := rw.(http.Hijacker).Hijack()
			defer conn.Close()

			if err := response.Write(readWriter); err != nil {
				log.Error("error writing response to remote: %v", err)
				return
			}

			if err := readWriter.Flush(); err != nil {
				log.Error("error flushing response to remote: %v", err)
				return
			}

			log.Debug("done")
		} else {
			rw.WriteHeader(http.StatusNotFound)
		}
	})

	return http.ListenAndServe("0.0.0.0:8080", nil)
}
