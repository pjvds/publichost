package client

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/pjvds/publichost/proto"
	"net"
	"sync"
	"sync/atomic"
)

var (
	address       = flag.String("address", "localhost:80", "The local address to make publicly available")
	remoteAddress = "publichost.me:8080"
)

type Tunnel struct {
	// The identifier of the tunnel which
	// is unique within a single client context.
	Id proto.TunnelId

	// The public address. This is available
	// after the tunnel is confirmed.
	PublicAddress string

	// The local address. This is given
	// at construct.
	LocalAddress string
}

func NewTunnel(id proto.TunnelId, localAddress string) *Tunnel {
	return &Tunnel{
		Id:           id,
		LocalAddress: localAddress,
	}
}

type PublicHostClient struct {
	conn net.Conn

	nextTunnelId int32

	reader *proto.Reader
	writer *proto.Writer

	outgoing chan proto.Envelop

	tunnelsLock sync.RWMutex
	tunnels     map[proto.TunnelId]*Tunnel
}

func Dial(address string) (*PublicHostClient, error) {
	resolved, err := net.ResolveTCPAddr("tcp", remoteAddress)
	if err != nil {
		return nil, fmt.Errorf("Unable to resolve %v: %v", remoteAddress, err)
	}

	conn, err := net.DialTCP("tcp", nil, resolved)
	if err != nil {
		return nil, fmt.Errorf("Unable to connect publichost server: %v", err)
	}

	if err := handshake(conn); err != nil {
		return nil, err
	}

	bufferedReader := bufio.NewReader(conn)

	return &PublicHostClient{
		conn:   conn,
		reader: proto.NewReader(bufferedReader),
	}, nil
}

func handshake(conn net.Conn) error {
	// This is for future use. We don't have
	// anything to exchange during handshake.
	return nil
}

func (p *PublicHostClient) processIncomming() error {
	for {
		message, err := p.reader.Read()
		if err != nil {
			return err
		}

		p.dispachIncomming(message)
	}
}

func (p *PublicHostClient) dispachIncomming(envelop *proto.Envelop) {
	p.tunnelsLock.Lock()
	defer p.tunnelsLock.Unlock()

	if t, ok := p.tunnels[envelop.Header.TunnelId]; ok {
		// TODO:
		fmt.Printf("< MSG %v\n", envelop.Header.TunnelId)
		//t.incomming <- message.Payload
	} else {
		// TODO:
		// else we just discard the message and the caller
		// will never receive a reply.
	}
}

func (p *PublicHostClient) writeOutgoing() error {
	for {
		message := <-p.outgoing

		if err := p.writer.Write(message); err != nil {
			return err
		}
	}
}

func (p *PublicHostClient) OpenTunnel(hostname string, port int) *Tunnel {
	request := proto.NewExposeRequest(fmt.Sprintf("%v:%v", hostname, port))

	tunnelId := atomic.AddInt32(&p.nextTunnelId, 1)
	tunnel := NewTunnel(proto.TunnelId(tunnelId), fmt.Sprintf("%v:%v", hostname, port))

	// TODO: Lock
	p.tunnels[tunnel.Id] = tunnel

	return tunnel
}
