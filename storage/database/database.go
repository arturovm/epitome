package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3" // imported for driver only
	"github.com/pkg/errors"
	"github.com/pressly/goose"
	log "github.com/sirupsen/logrus"
)

const driverName = "sqlite3"

// Migrate attempts to apply the migrations in the given directory to the
// database up to the given version.
func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect(driverName)
	if err != nil {
		return errors.Wrap(err, "error setting goose dialect")
	}
	goose.SetLogger(log.StandardLogger())
	return goose.Up(db, dir)
}
