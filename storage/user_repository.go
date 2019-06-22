package storage

import (
	"github.com/google/uuid"

	"github.com/arturovm/epitome"
)

type UserRepository interface {
	Add(epitome.User) error
	ByID(uuid.UUID) (*epitome.User, error)
}
