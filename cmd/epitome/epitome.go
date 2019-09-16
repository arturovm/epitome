package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"

	"github.com/arturovm/epitome/conf"
	"github.com/arturovm/epitome/filemanager"
	"github.com/arturovm/epitome/server"
	"github.com/arturovm/epitome/storage/database"
	"github.com/arturovm/epitome/users"
)

func main() {
	// check if help flag is set
	if conf.Help {
		conf.PrintHelp()
		os.Exit(0)
	}
	// check if debug mode is enabled
	if conf.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug mode enabled")
	}

	// setup data dir
	dataDir := conf.DataDir()

	log.WithField("path", dataDir).Debug("initializing data directory")
	err := filemanager.TouchDir(dataDir)
	if err != nil {
		log.WithError(err).Fatal("failed to initialize data directory")
	}

	// setup storage manager
	dbFilename := "file:" + filepath.Join(dataDir, "data.db")

	log.WithField("path", dbFilename).Debug("connecting to database")
	db, err := sql.Open("sqlite3", dbFilename)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to database")
	}

	// run migrations
	migrationsDir := conf.MigrationsDir()

	log.WithField("path", migrationsDir).Debug("running migrations")
	err = database.Migrate(db, migrationsDir)
	if err != nil {
		log.WithError(err).Fatal("failed to run migrations")
	}

	// setup user repository
	userRepository := database.NewUserRepository(db)

	// setup users service
	usrs := users.New(userRepository)

	// setup users handler set
	_ = server.NewUsersHandlerSet(usrs)

	// start server
	log.WithFields(log.Fields{
		"address": conf.Addr,
		"port":    conf.Port,
	}).Info("server starting")

	addr := net.JoinHostPort(conf.Addr, fmt.Sprintf("%d", conf.Port))
	log.WithError(http.ListenAndServe(addr, nil)).
		Fatal("server error")
}
