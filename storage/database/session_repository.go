package database

import (
	"github.com/Masterminds/squirrel"
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
