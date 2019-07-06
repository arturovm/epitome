package server

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type requestUser struct {
	Username string
	Password string
}

func (h *APIHandler) postUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dec := json.NewDecoder(r.Body)

	var user requestUser
	err := dec.Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = h.api.Users().SignUp(user.Username, user.Password)
	if err != nil {
		log.WithField("error", err).Error("error creating user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *APIHandler) getUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")
	user, err := h.api.Users().UserInfo(username)
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
