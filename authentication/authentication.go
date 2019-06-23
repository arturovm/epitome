package authentication

import (
	"github.com/pkg/errors"

	"github.com/arturovm/epitome"
	"github.com/arturovm/epitome/storage"
)

var ErrWrongCredentials = errors.New("wrong username or password")

type Authentication struct {
	sessions storage.SessionRepository
	users    storage.UserRepository
}

func New(sessions storage.SessionRepository, users storage.UserRepository) *Authentication {
	return &Authentication{
		sessions: sessions,
		users:    users,
	}
}

func (a *Authentication) LogIn(username, password string) (*epitome.Session, error) {
	user, err := a.users.ByUsername(username)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving user")
	}

	if !user.PasswordMatch(password) {
		return nil, ErrWrongCredentials
	}

	session, err := epitome.NewSession(username)
	if err != nil {
		return nil, errors.Wrap(err, "error creating session")
	}

	err = a.sessions.Add(*session)
	if err != nil {
		return nil, errors.Wrap(err, "error saving session")
	}

	return session, nil
}
