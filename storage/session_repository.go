package storage

import "github.com/arturovm/epitome"

// SessionRepository represents a repository that handles sessions.
type SessionRepository interface {
	Add(epitome.Session) error
}
