package server

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"net"
)

type EchoRoundtrip struct {
	conn net.Conn
}

func (e *EchoRoundtrip) Serve() (n int64, err error) {
	defer e.conn.Close()

	if n, err = io.Copy(e.conn, e.conn); err != nil {
		log.Error("error echo'ing: %v", err)
	}

	return
}

type NetEchoService struct {
	listener net.Listener
}

func (n *NetEchoService) Addr() net.Addr {
	return n.listener.Addr()
}

func NewNetEchoServiceTCP(address string) (service *NetEchoService, err error) {
	var listener net.Listener
	if listener, err = net.Listen("tcp", address); err != nil {
		return
	}
	service = &NetEchoService{
		listener: listener,
	}
	return
}

func (n *NetEchoService) ServeSingle() (int64, error) {
	defer n.Close()

	roundtrip, err := n.accept()
	if err != nil {
		return 0, err
	}

	return roundtrip.Serve()
}

func (n *NetEchoService) Serve() (err error) {
	defer n.Close()

	var roundtrip *EchoRoundtrip

	for {
		if roundtrip, err = n.accept(); err != nil {
			return
		}
		go roundtrip.Serve()
	}
}

func (n *NetEchoService) Close() {
	if err := n.listener.Close(); err != nil {
		log.Error("error closing listener: %v", err)
	}
}

func (n *NetEchoService) accept() (roundtrip *EchoRoundtrip, err error) {
	var conn net.Conn
	if conn, err = n.listener.Accept(); err != nil {
		return
	}

	roundtrip = &EchoRoundtrip{
		conn: conn,
	}
	return
}

var _ = Describe("Send data to stream", func() {
	var manager *StreamManager
	var echoService *NetEchoService

	BeforeEach(func() {
		var err error
		dialer := newNetworkStreamDialer()
		manager = NewStreamManager(dialer)

		echoService, err = NewNetEchoServiceTCP("127.0.0.1:0")
		Expect(err).ToNot(HaveOccured())

		go func() {
			if err := echoService.Serve(); err != nil {
				Expect(err).ToNot(HaveOccured())
			}
		}()
	})
	AfterEach(func() {
		echoService.Close()
	})

	It("shoud send data", func() {
		streamId, err := manager.OpenStream(echoService.Addr().Network(), echoService.Addr().String())

		Expect(streamId).ToNot(Equal(uint32(0)))
		Expect(err).ToNot(HaveOccured())
	})
})
