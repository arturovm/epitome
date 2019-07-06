package users_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mockstorage "github.com/arturovm/epitome/storage/mock"
	"github.com/arturovm/epitome/users"
)

func TestNewUsersService(t *testing.T) {
	var repo mockstorage.UserRepository
	us, err := users.New(&repo)
	require.NoError(t, err)
	require.NotNil(t, us)
	repo.AssertExpectations(t)
}

func TestNewUsersServiceNoRepo(t *testing.T) {
	us, err := users.New(nil)
	require.Error(t, err)
	require.Nil(t, us)
}

func TestSignUp(t *testing.T) {
	var repo mockstorage.UserRepository
	repo.On("Add", mock.AnythingOfType("epitome.User")).Return(nil)

	us, _ := users.New(&repo)

	var username, password = "username", "password"
	user, err := us.SignUp(username, password)
	require.NoError(t, err)
	require.NotNil(t, user)

	require.NotNil(t, user.Credentials())

	repo.AssertExpectations(t)
}

/*
To-do list:

- Sign up when username already exists
- Sign up with invalid password
*/
