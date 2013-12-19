package server

import (
	"net/http"
)

type TunnelServer struct {
}

type Server interface {
	ListenAndServe() error
}

type server struct {
	address string // The TCP address to listen on,
}

func (s *server) ListenAndServe() error {
	return http.ListenAndServe(s.address, s)
}

func (s *server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.URL.Host != "tunneler.publichost.me" {
		response.WriteHeader(http.StatusNotFound)
		return
	}

	var hostname string
	if hostname = request.Header.Get("hostname"); hostname == "" {
		response.WriteHeader(http.StatusBadRequest)
		return
	}

	response.Header().Set("hostname", hostname)
}
