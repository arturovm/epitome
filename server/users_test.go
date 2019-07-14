package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arturovm/epitome"
	"github.com/arturovm/epitome/server"
	"github.com/arturovm/epitome/users/mock"
)

func TestNewUsersHandlerSet(t *testing.T) {
	usrs := new(mock.Users)
	handlerSet := server.NewUsersHandlerSet(usrs)
	require.NotNil(t, handlerSet)
	usrs.AssertExpectations(t)
}

func TestEmptyUserSignUp(t *testing.T) {
	var (
		usrs = new(mock.Users)

		url = "http://localhost:8080/api/users"
		w   = httptest.NewRecorder()

		buf = bytes.NewBuffer(nil)
	)
	handlerSet := server.NewUsersHandlerSet(usrs)

	req := httptest.NewRequest("POST", url, buf)

	handlerSet.SignUp(w, req)
	require.Equal(t, w.Code, http.StatusBadRequest)
	usrs.AssertExpectations(t)
}

func TestSignUp(t *testing.T) {
	var (
		usrs = new(mock.Users)

		url = "http://localhost:8080/api/users"
		w   = httptest.NewRecorder()

		buf     = bytes.NewBuffer(nil)
		reqUser = server.RequestUser{
			Username: "testusername",
			Password: "testpassword",
		}
		respUser = &epitome.User{Username: "testusername"}
	)
	usrs.On("SignUp", reqUser.Username, reqUser.Password).
		Return(respUser, nil)
	handlerSet := server.NewUsersHandlerSet(usrs)

	_ = json.NewEncoder(buf).Encode(reqUser)
	req := httptest.NewRequest("POST", url, buf)

	handlerSet.SignUp(w, req)
	require.Equal(t, w.Code, http.StatusCreated)
	usrs.AssertExpectations(t)

	var resp epitome.User
	err := json.NewDecoder(w.Body).Decode(&resp)
	require.NoError(t, err)
	require.Equal(t, resp.Username, reqUser.Username)
}

/*
To-do list:

- Conflicting user sign-up
*/
