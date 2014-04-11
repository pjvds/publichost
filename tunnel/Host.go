package tunnel

import ()

type Host interface {
    OpenTunnel(address string) (hostname string, err error)
	Serve() error
}
