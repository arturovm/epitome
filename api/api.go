package api

import (
	"github.com/arturovm/epitome/authentication"
	"github.com/arturovm/epitome/storage/database"
	"github.com/arturovm/epitome/users"
	"github.com/pkg/errors"
)

// API represents the API port.
type API struct {
	users          *users.Users
	authentication *authentication.Authentication
}

// New takes a database manager and returns an initialized API.
func New(m *database.Manager) (*API, error) {
	usrs, err := users.New(m.UserRepository)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize users service")
	}

	auth, err := authentication.New(m.SessionRepository, m.UserRepository)
	if err != nil {
		return nil, errors.Wrap(err,
			"cannot intialize authentication service")
	}
	return &API{
		users:          usrs,
		authentication: auth,
	}, nil
}

// Users returns the users service.
func (a *API) Users() *users.Users {
	return a.users
}

func (a *API) Authentication() *authentication.Authentication {
	return a.authentication
}
