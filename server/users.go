package server

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"github.com/arturovm/epitome/users"
)

type usersHandler struct {
	users *users.Users
}

type requestUser struct {
	Username string
	Password string
}

func (u *usersHandler) postUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dec := json.NewDecoder(r.Body)

	var user requestUser
	err := dec.Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = u.users.SignUp(user.Username, user.Password)
	if err != nil {
		log.WithField("error", err).Error("error creating user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (u *usersHandler) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")
	user, err := u.users.UserInfo(username)
	if err != nil {
		log.WithField("error", err).Error("error getting user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	enc := json.NewEncoder(w)
	err = enc.Encode(user)
	if err != nil {
		log.WithField("error", err).Error("error encoding user")
		return
	}
}
