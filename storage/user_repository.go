package storage

import (
	"github.com/arturovm/epitome"
)

// UserRepository represents a repository that handles users.
type UserRepository interface {
	Add(epitome.User) error
	ByUsername(string) (*epitome.User, error)
}
