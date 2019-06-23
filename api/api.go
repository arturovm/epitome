package api

import (
	"github.com/arturovm/epitome/storage/database"
	"github.com/arturovm/epitome/users"
)

// API represents the API port.
type API struct {
	users *users.Users
}

// New takes a database manager and returns an initialized API.
func New(m *database.Manager) *API {
	return &API{
		users: users.New(m.UserRepository),
	}
}

// Users returns the users service.
func (a *API) Users() *users.Users {
	return a.users
}
