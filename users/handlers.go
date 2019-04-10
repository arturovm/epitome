package users

import (
	"encoding/json"
	"net/http"

	"github.com/arturovm/epitome/data/user"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

// PostUser handles POST requests to the users endpoint
func PostUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// decode request body
	var ru requestUser
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&ru)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create user in database
	u, err = user.Create(ru.Username, ru.Password)
	if err != nil {
		switch err {
		case user.ErrInvalidPassword:
			w.WriteHeader(http.StatusBadRequest)
		case user.ErrPasswordHashingFailed:
			w.WriteHeader(http.StatusInternalServerError)
		case user.ErrUserExists:
			w.WriteHeader(http.StatusConflict)
		default:
			log.WithField("error", err).Error("error creating user")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)

	// encode user in response
	enc := json.NewEncoder(w)
	err = enc.Encode(u)
	if err != nil {
		log.WithField("error", err).Error("error encoding user")
	}

	return
}
