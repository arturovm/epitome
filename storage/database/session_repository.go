package database

import (
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/arturovm/epitome"
)

type SessionRepository struct {
	db *sqlx.DB
}

func (r *SessionRepository) Add(session epitome.Session) error {
	_, err := squirrel.Insert("sessions").
		Columns("id", "key", "username").
		Values(session.ID, session.Key, session.Username).
		RunWith(r.db).
		Exec()
	return err
}
func (r *SessionRepository) ByID(id uuid.UUID) (*epitome.Session, error) {
	return nil, nil
}
func (r *SessionRepository) ByUsername(username string) ([]*epitome.Session, error) {
	return nil, nil
}
