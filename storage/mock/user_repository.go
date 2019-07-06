package mock

import (
	"github.com/stretchr/testify/mock"

	"github.com/arturovm/epitome"
)

// UserRepository is a mock implementation of storage.UserRepository.
type UserRepository struct {
	mock.Mock
}

// Add implements UserRepository.Add.
func (r *UserRepository) Add(user epitome.User) error {
	args := r.Mock.Called(user)
	return args.Error(0)
}

// ByUsername implements UserRepository.ByUsername.
func (r *UserRepository) ByUsername(username string) (*epitome.User, error) {
	args := r.Mock.Called(username)
	return args.Get(0).(*epitome.User), args.Error(1)
}
