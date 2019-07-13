package epitome

import (
	"strings"
)

// User represents a system user.
type User struct {
	Username    string
	credentials *Credentials
}

// CreateUser takes a username and password and returns a new user with valid
// credentials.
func CreateUser(username, password string) (*User, error) {
	creds, err := GenerateCredentials(password)
	if err != nil {
		return nil, err
	}

	u := NewUser(username, creds)
	return &u, nil
}

// NewUser takes a username and credentials and returns a user with credentials
// set to the specified credentials.
func NewUser(username string, credentials *Credentials) User {
	return User{
		Username:    strings.ToLower(username),
		credentials: credentials,
	}
}

// Credentials returns a user's credentials.
func (u *User) Credentials() *Credentials {
	return u.credentials
}
