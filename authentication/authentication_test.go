package authentication_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/arturovm/epitome"
	"github.com/arturovm/epitome/authentication"
	"github.com/arturovm/epitome/storage"
	mockstorage "github.com/arturovm/epitome/storage/mock"
)

type AuthenticationTestSuite struct {
	suite.Suite
	users    *mockstorage.UserRepository
	sessions *mockstorage.SessionRepository
	auth     *authentication.Authentication
	username string
	password string
}

func (s *AuthenticationTestSuite) SetupTest() {
	s.users = new(mockstorage.UserRepository)
	s.sessions = new(mockstorage.SessionRepository)
	s.auth = authentication.New(s.sessions, s.users)
	s.username = "testuser"
	s.password = "testpassword"
}

func (s *AuthenticationTestSuite) TestLogInNoUser() {
	var emptyUser *epitome.User
	s.users.On("ByUsername", s.username).
		Return(emptyUser, storage.ErrUserNotFound)

	sess, err := s.auth.LogIn(s.username, s.password)

	require.EqualError(s.T(), err,
		authentication.ErrInvalidCredentials.Error())
	require.Nil(s.T(), sess)
	s.users.AssertExpectations(s.T())
	s.sessions.AssertExpectations(s.T())
}

func (s *AuthenticationTestSuite) TestLogIn(t *testing.T) {
	u, _ := epitome.CreateUser(s.username, s.password)

	s.users.On("ByUsername", u.Username).Return(u, nil)
	s.sessions.On("Add", mock.AnythingOfType("epitome.Session")).
		Return(nil)

	sess, err := s.auth.LogIn(s.username, s.password)

	require.NoError(s.T(), err)
	require.NotNil(s.T(), sess)
	require.Equal(s.T(), sess.Username, s.username)
	s.users.AssertExpectations(s.T())
	s.sessions.AssertExpectations(s.T())
}
func (s *AuthenticationTestSuite) TestLogInInvalidPassword(t *testing.T) {
	u, _ := epitome.CreateUser(s.username, s.password)

	s.users.On("ByUsername", u.Username).Return(u, nil)

	sess, err := s.auth.LogIn(s.username, "wrong password")

	require.EqualError(s.T(), err,
		authentication.ErrInvalidCredentials.Error())
	require.Nil(s.T(), sess)
	s.users.AssertExpectations(s.T())
	s.sessions.AssertExpectations(s.T())
}

func TestAuthentication(t *testing.T) {
	suite.Run(t, new(AuthenticationTestSuite))
}

/*
To-do list:

- Genuine errors retrieving a user should be reported as something different
  than ErrInvalidCredentials
*/
