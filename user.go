package epitome

import (
	"strings"
)

// User represents a system user.
type User struct {
	Username    string
	credentials *Credentials
}

// NewUser initializes a new user with the given parameters.
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
