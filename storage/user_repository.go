package storage

import (
	"github.com/arturovm/epitome"
)

type UserRepository interface {
	Add(epitome.User) error
	ByUsername(string) (*epitome.User, error)
}
