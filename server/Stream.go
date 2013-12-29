package server

type DataContext struct {
	Data []byte
}

type Stream struct {
	conn StreamConnection

    isClosed bool

	Outgoing  chan *DataContext
	Incomming chan *DataContext

	Id uint32
}

func newStream(id uint32, conn StreamConnection) *Stream {
	return &Stream{
		Id:   id,
		conn: conn,
	}
}

func (s *Stream) Serve() error {
	defer s.Close()

	e := make(chan error)

	go func() {
        if err := s.readIncomming(); err != nil {
            e <- err
        }
    } 

    go func() {
        if err := s.writeOutgoing(); err != nil {
            e <- err
        }
    }

    err := <- e
    return err
}

func (s *Stream) readIncomming() (err error) {
	var n int
	var buffer [2048]byte

	for !s.isClosed {
		if n, err := s.conn.Read(buffer); err != nil {
			log.Debug("error reading incomming stream data: %v", err)
			break
		}

		s.Incomming <- &DataContext{
			Data: buffer[0:n],
		}
	}

    if s.isClosed {
        log.Debug("cleared error because stream is closed")
        err = nil
    }

	return
}

func (s *Stream) Close() error {
	// TODO: Notificy owning connection that we are closed.
	return s.conn.Close()
}
