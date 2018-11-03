package data

import (
	"os"
	"path/filepath"

	"github.com/arturovm/epitome/conf"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const dbVersion int64 = 1

// Setup initializes necessary directories and files for data persistence
func Setup() error {
	dataDir := conf.DataDir()
	path := os.ExpandEnv(dataDir)
	log.WithField("path", path).Debug("initializing data directory")

	// touch data directory
	err := os.MkdirAll(path, os.ModeDir|0755)
	if err != nil {
		return errors.Wrap(err, "error creating data directory")
	}

	// initialize database "connection"
	filename := "file:" + filepath.Join(path, "data.db")
	err = openDB(filename)
	if err != nil {
		return errors.Wrap(err, "error opening database file")
	}

	// perform migrations
	err = migrate(dbVersion, conf.MigrationsDir())
	if err != nil {
		return errors.Wrap(err, "failed to run database migrations")
	}

	return nil
}
