package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/arturovm/epitome/api"
)

// APIHandler represents the API port HTTP adapter.
type APIHandler struct {
	api    *api.API
	router *httprouter.Router
}

// NewHandler takes an instance of the API port and returns an initialized
// handler.
func NewAPIHandler(api *api.API) *APIHandler {
	h := &APIHandler{
		api: api,
	}

	router := httprouter.New()
	h.registerRoutes(router)

	return h
}

func (h *APIHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.router.ServeHTTP(w, req)
}
