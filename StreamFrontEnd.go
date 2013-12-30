package publichost

import (
	"github.com/op/go-logging"
	"io"
	"net"
)

type StreamId uint32

type StreamData struct {
	Data []byte

	EOF bool

	Ack chan bool
	Err chan error
}

func NewStreamData(data []byte, eof bool) *StreamData {
	return &StreamData{
		Data: data,
		EOF:  eof,
		Ack:  make(chan bool, 1),
		Err:  make(chan error, 1),
	}
}

type StreamFrondEnd struct {
	Id   StreamId
	conn *net.TCPConn

	log *logging.Logger

	Outgoing  chan *StreamData
	Incomming chan *StreamData

	outgoingDone chan bool
}

func Dial(id StreamId, laddr string) (stream *StreamFrondEnd, err error) {
	var conn *net.TCPConn
	var address *net.TCPAddr

	if address, err = net.ResolveTCPAddr("tcp", laddr); err != nil {
		return
	}

	if conn, err = net.DialTCP("tcp", nil, address); err != nil {
		return
	}

	stream = &StreamFrondEnd{
		Id:   id,
		conn: conn,
		log:  log,

		Outgoing:  make(chan *StreamData, 25),
		Incomming: make(chan *StreamData, 25),

		outgoingDone: make(chan bool, 1),
	}
	return
}

// Serves incomming and outgoing data until both of them
// reached EOF or errored.
func (s *StreamFrondEnd) Serve() {
	readClosed := make(chan bool, 1)
	go func() {
		if err := s.read(); err != nil {
			log.Debug("reading stopped with error: %v", err)
		} else {
			log.Debug("reading finished")
		}

		readClosed <- true
	}()

	for {
		select {
		case outgoing := <-s.Outgoing:
			s.handleOutgoing(outgoing)
		case <-s.outgoingDone:
			break
		}
	}

	// Wait until we finished reading
	<-readClosed
}

func (s *StreamFrondEnd) StreamData(data []byte) (ack chan bool, err chan error) {
	d := NewStreamData(data, false)
	s.Outgoing <- d

	return d.Ack, d.Err
}

func (s *StreamFrondEnd) handleOutgoing(outgoing *StreamData) {
	if err := s.write(outgoing.Data); err != nil {
		outgoing.Err <- err

		s.log.Debug("Closing writer because of error: %v", err)
		s.closeWrite()
	} else if outgoing.EOF {
		s.log.Debug("Closing writer because of EOF")
		s.closeWrite()
	}
}

func (s *StreamFrondEnd) write(data []byte) error {
	written, err := s.conn.Write(data)
	if err != nil {
		s.log.Debug("writing error: ", err)
	} else {
		s.log.Debug("%v bytes written to %v", written, s.conn.RemoteAddr())
	}

	return err
}

func (s *StreamFrondEnd) read() (err error) {
	defer s.closeRead()

	var read int

	buffer := make([]byte, 8*1024)
	eof := false

	for {
		if read, err = s.conn.Read(buffer); err != nil {
			if err == io.EOF {
				eof = true
				err = nil
			}
		}

		data := buffer[0:read]
		s.Incomming <- NewStreamData(data, eof)
	}

}

func (s *StreamFrondEnd) closeWrite() {
	if err := s.conn.CloseWrite(); err != nil {
		s.log.Debug("close write error: %v", err)
	}

	close(s.Outgoing)
	s.Outgoing = nil

	s.outgoingDone <- true
	s.log.Debug("closed write")
}

func (s *StreamFrondEnd) closeRead() {
	if err := s.conn.CloseWrite(); err != nil {
		s.log.Debug("close read error: %v", err)
	}

	close(s.Incomming)
}

func (s *StreamFrondEnd) Close() {
	s.conn.Close()
}
