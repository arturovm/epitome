package storage

import (
	"github.com/google/uuid"

	"github.com/arturovm/epitome"
)

// SessionRepository represents a repository that handles sessions.
type SessionRepository interface {
	Add(epitome.Session) error
	ByID(uuid.UUID) (*epitome.Session, error)
	ByUsername(string) ([]*epitome.Session, error)
}
