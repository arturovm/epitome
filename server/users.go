package server

import (
	"encoding/json"
	"net/http"

	"github.com/arturovm/epitome/storage"
	"github.com/arturovm/epitome/users"
)

type UsersHandlerSet struct {
	users users.Service
}

type RequestUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUsersHandlerSet(users users.Service) *UsersHandlerSet {
	return &UsersHandlerSet{users: users}
}

func (h *UsersHandlerSet) SignUp(w http.ResponseWriter, req *http.Request) {
	var reqUser RequestUser
	err := json.NewDecoder(req.Body).Decode(&reqUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.users.SignUp(reqUser.Username, reqUser.Password)
	if err == storage.ErrUserExists {
		w.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}
