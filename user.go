package epitome

import (
	"strings"
)

// User represents a system user.
type User struct {
	Username    string
	credentials *Credentials
}

func CreateUser(username, password string) (*User, error) {
	creds, err := GenerateCredentials(password)
	if err != nil {
		return nil, err
	}

	u := NewUser(username, creds)
	return &u, nil
}

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
