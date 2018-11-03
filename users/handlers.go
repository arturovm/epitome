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
	_, err = user.Create(ru.Username, ru.Password)
	if err != nil {
		switch err {
		case user.ErrInvalidPassword:
			w.WriteHeader(http.StatusBadRequest)
			return
		case user.ErrPasswordHashingFailed:
			w.WriteHeader(http.StatusInternalServerError)
			return
		case user.ErrUserExists:
			w.WriteHeader(http.StatusConflict)
			return
		default:
			log.WithField("error", err).Error("error creating user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	return
}
