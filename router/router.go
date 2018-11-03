package router

import (
	"github.com/arturovm/epitome/auth"
	"github.com/arturovm/epitome/users"
	"github.com/julienschmidt/httprouter"
)

// Get returns a router configured and ready to use
func Get() *httprouter.Router {
	r := httprouter.New()

	// auth
	r.POST("/api/auth/sessions", auth.PostSession)
	r.DELETE("/api/auth/sessions/:token", auth.DeleteSession)
	// users
	r.POST("/api/users", users.PostUser)

	return r
}
