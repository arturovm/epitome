package server

import (
	"net"
	"net/http"
	"strconv"

	"github.com/arturovm/epitome/api"
	"github.com/julienschmidt/httprouter"
)

// Server represents the server adapter for the API port.
type Server struct {
	router *httprouter.Router
}

// New takes an API instance and returns an initialized server.
func New(api *api.API) *Server {
	router := setupRoutes(api)
	return &Server{
		router: router,
	}
}

// Start takes a host and a port and starts an HTTP server on the resulting
// address.
func (s *Server) Start(host string, port int) error {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	return http.ListenAndServe(addr, s.router)
}
