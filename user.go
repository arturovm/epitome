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
func NewUser(username string) User {
	return User{
		Username: strings.ToLower(username),
	}
}

// Credentials returns a user's credentials.
func (u *User) Credentials() *Credentials {
	return u.credentials
}

// SetCredentials sets a user's credentials to the supplied value.
func (u *User) SetCredentials(credentials *Credentials) {
	u.credentials = credentials
}
