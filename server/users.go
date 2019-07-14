package server

import (
	"encoding/json"
	"net/http"

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

	user, _ := h.users.SignUp(reqUser.Username, reqUser.Password)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}
