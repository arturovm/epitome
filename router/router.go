package router

import (
	"github.com/ArturoVM/epitome/auth"
	"github.com/julienschmidt/httprouter"
)

func Get() *httprouter.Router {
	r := httprouter.New()

	r.POST("/api/auth/sessions", auth.PostSessions)
	r.DELETE("/api/auth/sessions/:token", auth.DeleteSession)

	return r
}
