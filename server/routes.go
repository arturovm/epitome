package server

import "github.com/julienschmidt/httprouter"

func (h *APIHandler) registerRoutes(router *httprouter.Router) {
	router.POST("/api/users", h.postUser)
	router.GET("/api/users/:username", h.getUser)
}
