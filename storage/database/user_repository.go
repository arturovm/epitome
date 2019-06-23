package database

import (
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/arturovm/epitome"
)

type UserRepository struct {
	db *sqlx.DB
}

func (r *UserRepository) Add(user epitome.User) error {
	_, err := squirrel.Insert("users").
		Columns("username", "password", "salt").
		Values(user.Username, user.Password, user.Salt).
		RunWith(r.db).
		Exec()
	return err
}

func (r *UserRepository) ByUsername(username string) (*epitome.User, error) {
	query, args, err := squirrel.Select("username", "password", "salt").
		From("users").
		Where(squirrel.Eq{"username": username}).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "error building select query")
	}

	var user epitome.User
	err = r.db.Get(&user, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "error querying database")
	}

	return &user, nil
}
