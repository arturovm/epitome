package server

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type loginRequest struct {
	Username string
	Password string
}

func (h *APIHandler) postSession(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var login loginRequest
	err := dec.Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session, err := h.api.
		Authentication().
		LogIn(login.Username, login.Password)
	if err != nil {
		log.WithField("error", err).Error("error logging user in")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusCreated)

	enc := json.NewEncoder(w)
	err = enc.Encode(session)
	if err != nil {
		log.WithField("error", err).Error("error encoding session")
		return
	}
}
