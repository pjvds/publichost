package tunnel

import (
	"errors"
	"github.com/pjvds/publichost/net/message"
	"io"
	"fmt"
	"net/http"
	"net"
	"bufio"
	"time"
)

type backendHost struct {
	conn   io.ReadWriteCloser
	reader message.Reader
	writer message.Writer

	tunnel       Tunnel
	LocalAddress string

	handlers map[byte]MessageHandler
}

func connect(address string) (conn net.Conn, err error) {
	var req *http.Request
	var response *http.Response

	if req, err = http.NewRequest("CONNECT", "/", nil); err != nil {
		return
	}

	req.Header.Add("X-PUBLICHOST", "true")
	req.Header.Add("X-PUBLICHOST-VERSION", "0.1")
	req.Header.Add("Keep-Alive", "")

	if conn, err = net.Dial("tcp4", address); err != nil {
		return
	}

	var resolved *net.TCPAddr
	var tcpConn *net.TCPConn
	if resolved, err = net.ResolveTCPAddr("tcp", address); err != nil {
		return
	}
	if tcpConn, err = net.DialTCP("tcp", nil, resolved); err != nil {
		return
	}
	if err = tcpConn.SetLinger(5); err != nil {
		return
	}
	conn = tcpConn

	defer func() {
		if err != nil {
			conn.Close()
		}
	}()

	bufRW := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
	if err = req.Write(bufRW); err != nil {
		return
	}
	if err = bufRW.Writer.Flush(); err != nil {
		return
	}
	conn.SetWriteDeadline(time.Time{})


	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	if response, err = http.ReadResponse(bufRW.Reader, req); err != nil {
		return
	}
	conn.SetReadDeadline(time.Time{})

	if response.StatusCode != 200 {
		err = fmt.Errorf("Unexpected response: %v", response.StatusCode)
		return
	}
	return
}

func NewBackendHost(address string) (host Host, err error) {
	var conn net.Conn
	if conn, err = connect(address); err != nil {
		return
	}

	h := &backendHost{
		conn:     conn,
		reader:   message.NewReader(conn),
		writer:   message.NewWriter(conn),
		tunnel:   NewTunnelBackend(),
		handlers: make(map[byte]MessageHandler),
	}
	h.handlers[message.OpOpenStream] = NewOpenStreamHandler(h.tunnel)
	h.handlers[message.OpCloseStream] = NewCloseStreamHandler(h.tunnel)
	h.handlers[message.OpWriteStream] = NewWriteStreamHandler(h.tunnel)
	h.handlers[message.OpReadStream] = NewReadStreamHandler(h.tunnel)

	host = h
	return
}

func (h *backendHost) OpenTunnel(localAddress string) (hostname string, err error) {
	var response *message.Message

	request := message.NewMessage(message.OpOpenTunnel, 1, []byte(localAddress))
	if err = h.writer.Write(request); err != nil {
		log.Debug("unable to write handshake message: %v", err)
		return
	}

	if response, err = h.reader.Read(); err != nil {
		log.Debug("unable to read response to handshake: %v", err)
		return
	}

	if response.TypeId == message.Nack {
		err = errors.New(string(response.Body))
		log.Debug("unable to open tunnel: %v", err)
		return
	}

	hostname = string(response.Body)

	log.Debug("tunnel opened successfully")
	return
}

func (h *backendHost) Serve() (err error) {
	var request *message.Message

	for {
		if request, err = h.reader.Read(); err != nil {
			break
		}

		fmt.Printf("incomming: %s\n", request)

		response := NewResponseWriter(h.writer, request.CorrelationId)

		if handler, ok := h.handlers[request.TypeId]; ok {
			go handler.Handle(response, request)
		} else {
			log.Debug("No handlers for message type id %v", request.TypeId)
			go response.Nack(errors.New("unknown type id"))
		}
	}

	return
}
