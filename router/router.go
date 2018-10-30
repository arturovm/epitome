package router

import (
	"github.com/arturovm/epitome/auth"
	"github.com/julienschmidt/httprouter"
)

// Get returns a router configured and ready to use
func Get() *httprouter.Router {
	r := httprouter.New()

	r.POST("/api/auth/sessions", auth.PostSessions)
	r.DELETE("/api/auth/sessions/:token", auth.DeleteSession)

	return r
}
