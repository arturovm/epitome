package database

import (
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/arturovm/epitome"
)

type UserRepository struct {
	db *sqlx.DB
}

func (r *UserRepository) Add(user epitome.User) error {
	id := uuid.New()
	_, err := squirrel.Insert("users").
		Columns("id", "username", "password", "salt").
		Values(id, user.Username, user.Password, user.Salt).
		RunWith(r.db).
		Exec()
	return err
}
