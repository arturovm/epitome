package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/arturovm/epitome"
)

type UserRepository struct {
	mock.Mock
}

func (r *UserRepository) Add(user epitome.User) error {
	args := r.Mock.Called(user)
	return args.Error(0)
}

func (r *UserRepository) ByUsername(username string) (*epitome.User, error) {
	args := r.Mock.Called(username)
	return args.Get(0).(*epitome.User), args.Error(0)
}
