package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/arturovm/epitome"
)

type Users struct {
	mock.Mock
}

func (us *Users) SignUp(username, password string) (*epitome.User, error) {
	args := us.Called(username, password)
	return args.Get(0).(*epitome.User), args.Error(1)
}
