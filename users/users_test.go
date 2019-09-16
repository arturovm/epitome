package users_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/arturovm/epitome/storage"
	mockstorage "github.com/arturovm/epitome/storage/mock"
	"github.com/arturovm/epitome/users"
)

func TestSignUp(t *testing.T) {
	var repo mockstorage.UserRepository
	repo.On("Add", mock.AnythingOfType("epitome.User")).Return(nil)

	us := users.New(&repo)

	var username, password = "username", "password"
	user, err := us.SignUp(username, password)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.NotNil(t, user.Credentials())
	repo.AssertExpectations(t)
}

func TestSignUpUserExists(t *testing.T) {
	var repo mockstorage.UserRepository
	repo.On("Add", mock.AnythingOfType("epitome.User")).
		Return(storage.ErrUserExists)

	us := users.New(&repo)

	username := "existingUser"
	user, err := us.SignUp(username, "")
	require.EqualError(t, err, storage.ErrUserExists.Error())
	require.Nil(t, user)
	repo.AssertExpectations(t)
}

/*
To-do list:

- Sign up when username already exists
- Sign up with invalid password
*/
