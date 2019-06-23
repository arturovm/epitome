package server

import (
	"net"
	"net/http"
	"strconv"

	"github.com/arturovm/epitome/api"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	router *httprouter.Router
}

func New(api *api.API) *Server {
	router := setupRoutes(api)
	return &Server{
		router: router,
	}
}

func (s *Server) Start(host string, port int) error {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	return http.ListenAndServe(addr, s.router)
}
