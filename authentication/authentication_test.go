package authentication_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/arturovm/epitome"
	"github.com/arturovm/epitome/authentication"
	"github.com/arturovm/epitome/storage"
	mockstorage "github.com/arturovm/epitome/storage/mock"
)

func TestLogInNoUser(t *testing.T) {
	var (
		users    mockstorage.UserRepository
		sessions mockstorage.SessionRepository

		username = "notauser"
		password = "fakepassword"
	)

	var emptyUser *epitome.User
	users.On("ByUsername", username).
		Return(emptyUser, storage.ErrUserNotFound)

	auth := authentication.New(&sessions, &users)

	sess, err := auth.LogIn(username, password)
	require.EqualError(t, err,
		authentication.ErrInvalidCredentials.Error())
	require.Nil(t, sess)
	users.AssertExpectations(t)
	sessions.AssertExpectations(t)
}

func TestLogIn(t *testing.T) {
	var (
		users    mockstorage.UserRepository
		sessions mockstorage.SessionRepository

		username = "testuser"
		password = "testpassword"
	)

	u, _ := epitome.CreateUser(username, password)

	users.On("ByUsername", u.Username).Return(u, nil)
	sessions.On("Add", mock.AnythingOfType("epitome.Session")).
		Return(nil)

	auth := authentication.New(&sessions, &users)

	sess, err := auth.LogIn(username, password)
	require.NoError(t, err)
	require.NotNil(t, sess)
	require.Equal(t, sess.Username, username)
	users.AssertExpectations(t)
	sessions.AssertExpectations(t)
}
func TestLogInInvalidPassword(t *testing.T) {
	var (
		users    mockstorage.UserRepository
		sessions mockstorage.SessionRepository

		username = "testuser"
		password = "testpassword"
	)

	u, _ := epitome.CreateUser(username, password)

	users.On("ByUsername", u.Username).Return(u, nil)

	auth := authentication.New(&sessions, &users)

	sess, err := auth.LogIn(username, "wrong password")
	require.EqualError(t, err,
		authentication.ErrInvalidCredentials.Error())
	require.Nil(t, sess)
	users.AssertExpectations(t)
	sessions.AssertExpectations(t)
}

/*
To-do list:

- Genuine errors retrieving a user should be reported as something different
  than ErrInvalidCredentials
*/
