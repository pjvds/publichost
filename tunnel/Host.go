package tunnel

import ()

type Host interface {
	Serve() error
}
