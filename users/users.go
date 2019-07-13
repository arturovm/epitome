package users

import (
	"github.com/pkg/errors"

	"github.com/arturovm/epitome"
	"github.com/arturovm/epitome/storage"
)

// Users is a user managing service.
type Users struct {
	users storage.UserRepository
}

// ErrInvalidPassword is returned when the given password doesn't satisfy
// the minimum criteria.
var ErrInvalidPassword = errors.New("password is invalid")

// New takes a user repository and returns an initialized users service.
func New(users storage.UserRepository) *Users {
	return &Users{users: users}
}

// SignUp attempts to create a new user with the given username and password.
func (u *Users) SignUp(username, password string) (*epitome.User, error) {
	user, err := epitome.CreateUser(username, password)
	if err != nil {
		return nil, errors.Wrap(err, "error creating user")
	}

	err = u.users.Add(*user)
	if err != nil {
		return nil, errors.Wrap(err, "error saving user")
	}
	return user, nil
}
