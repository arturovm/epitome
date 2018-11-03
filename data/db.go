package data

import (
	"github.com/gocraft/dbr"
	_ "github.com/mattn/go-sqlite3" // import for driver only
	"github.com/pressly/goose"
	log "github.com/sirupsen/logrus"
)

var conn *dbr.Connection

func openDB(path string) error {
	var err error
	conn, err = dbr.Open("sqlite3", path, nil)
	if err != nil {
		return err
	}
	return nil
}

func migrate(v int64, dir string) error {
	goose.SetDialect("sqlite3")
	goose.SetLogger(log.StandardLogger())
	return goose.UpTo(conn.DB, dir, v)
}

// GetSession returns a database connection session. For non-atomic operations.
// If you need atomicity, use GetTx.
func GetSession() *dbr.Session {
	return conn.NewSession(nil)
}

// GetTx returns an initialized database transaction. Use this for atomic
// operations.
func GetTx() (*dbr.Tx, error) {
	return GetSession().Begin()
}
