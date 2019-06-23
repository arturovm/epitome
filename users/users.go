package users

import (
	"github.com/pkg/errors"

	"github.com/arturovm/epitome"
	"github.com/arturovm/epitome/storage"
)

// Users is a user managing service.
type Users struct {
	repository storage.UserRepository
}

// New takes a user repository and returns an initializes users service.
func New(repository storage.UserRepository) *Users {
	return &Users{repository: repository}
}

// SignUp attempts to create a new user with the given username and password.
func (u *Users) SignUp(username, password string) error {
	user, err := epitome.NewUser(username, password)
	if err != nil {
		return errors.Wrap(err, "error creating new user")
	}

	err = u.repository.Add(*user)
	if err != nil {
		return errors.Wrap(err, "error saving user")
	}
	return nil
}

// UserInfo retrieves a user instance from the database with the given username.
func (u *Users) UserInfo(username string) (*epitome.User, error) {
	return u.repository.ByUsername(username)
}
