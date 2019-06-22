package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // import for driver only
	"github.com/pressly/goose"
	log "github.com/sirupsen/logrus"
)

type Manager struct {
	*UserRepository
}

const driverName = "sqlite3"

func New(path string) (*Manager, error) {
	db, err := sqlx.Connect(driverName, path)
	if err != nil {
		return nil, err
	}
	return &Manager{&UserRepository{db}}, nil
}

func (m *Manager) Migrate(v int64, dir string) error {
	goose.SetDialect(driverName)
	goose.SetLogger(log.StandardLogger())
	return goose.UpTo(m.db.DB, dir, v)
}
