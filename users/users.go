package users

import (
	"github.com/pkg/errors"

	"github.com/arturovm/epitome"
	"github.com/arturovm/epitome/storage"
)

type Users struct {
	repository storage.UserRepository
}

func New(repository storage.UserRepository) *Users {
	return &Users{repository: repository}
}

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
