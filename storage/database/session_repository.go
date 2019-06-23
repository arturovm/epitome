package database

import (
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/arturovm/epitome"
)

// SessionRepository implements storage.SessionRepository.
type SessionRepository struct {
	db *sqlx.DB
}

// Add implements SessionRepository.Add.
func (r *SessionRepository) Add(session epitome.Session) error {
	_, err := squirrel.Insert("sessions").
		Columns("id", "key", "username").
		Values(session.ID, session.Key, session.Username).
		RunWith(r.db).
		Exec()
	return err
}

// ByID implements SessionRepository.ByID.
func (r *SessionRepository) ByID(id uuid.UUID) (*epitome.Session, error) {
	return nil, nil
}

// ByUsername implements SessionRepository.ByUsername.
func (r *SessionRepository) ByUsername(username string) ([]*epitome.Session, error) {
	return nil, nil
}
