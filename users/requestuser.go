package users

import (
	"github.com/arturovm/epitome/data/user"
)

// requestUser is the type used to read from an HTTP request
type requestUser struct {
	Username string
	Password string
}

// toUser initializes an user instance with properly formatted username and
// password parameters from a requestUser instance
func (ru *requestUser) toUser() (*user.User, error) {
	return user.New(ru.Username, ru.Password)
}
