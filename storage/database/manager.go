package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3" // import for driver only
	"github.com/pressly/goose"
	log "github.com/sirupsen/logrus"
)

// Manager is a storage manager that abstracts database details.
type Manager struct {
	db *sqlx.DB
	*UserRepository
	*SessionRepository
}

const driverName = "sqlite3"

// New takes a path and opens a sqlite3 connection to the given file.
func New(path string) (*Manager, error) {
	db, err := sqlx.Connect(driverName, path)
	if err != nil {
		return nil, err
	}
	return &Manager{
		db:                db,
		UserRepository:    &UserRepository{db},
		SessionRepository: &SessionRepository{db},
	}, nil
}

// Migrate attempts to apply the migrations in the given directory to the
// database up to the given version.
func (m *Manager) Migrate(v int64, dir string) error {
	goose.SetDialect(driverName)
	goose.SetLogger(log.StandardLogger())
	return goose.UpTo(m.db.DB, dir, v)
}
