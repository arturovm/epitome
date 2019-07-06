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

func (u *User) Credentials() *Credentials {
	return u.credentials
}

func (u *User) SetCredentials(credentials *Credentials) {
	u.credentials = credentials
}
