package server

import (
	"net"
	"net/http"
	"strconv"
)

// Server represents the server adapter for the API port.
type Server struct {
	handler *APIHandler
}

// New takes an API instance and returns an initialized server.
func New(handler *APIHandler) *Server {
	return &Server{
		handler: handler,
	}
}

// Start takes a host and a port and starts an HTTP server on the resulting
// address.
func (s *Server) Start(host string, port int) error {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	return http.ListenAndServe(addr, s.handler)
}
