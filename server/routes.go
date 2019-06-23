package server

import (
	"github.com/arturovm/epitome/api"
	"github.com/julienschmidt/httprouter"
)

func setupRoutes(api *api.API) *httprouter.Router {
	router := httprouter.New()

	uh := usersHandler{api.Users()}
	router.POST("/api/users", uh.postUser)
	router.GET("/api/users/:username", uh.getUser)

	return router
}
