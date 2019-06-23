package api

import (
	"github.com/arturovm/epitome/storage/database"
	"github.com/arturovm/epitome/users"
)

type API struct {
	users *users.Users
}

func New(m *database.Manager) *API {
	return &API{
		users: users.New(m.UserRepository),
	}
}

func (a *API) Users() *users.Users {
	return a.users
}
