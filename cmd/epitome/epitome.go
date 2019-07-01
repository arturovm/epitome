package main

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/arturovm/epitome/api"
	"github.com/arturovm/epitome/conf"
	"github.com/arturovm/epitome/filemanager"
	"github.com/arturovm/epitome/server"
	"github.com/arturovm/epitome/storage/database"
)

const dbVersion = 1

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
		log.WithField("error", err).
			Fatal("failed to initialize data directory")
	}

	// setup storage manager
	dbFilename := "file:" + filepath.Join(dataDir, "data.db")

	log.WithField("path", dbFilename).Debug("initializing storage manager")
	storageManager, err := database.New(dbFilename)
	if err != nil {
		log.WithField("error", err).
			Fatal("failed to initialize storage manager")
	}

	// run migrations
	migrationsDir := conf.MigrationsDir()

	log.WithField("path", migrationsDir).Debug("running migrations")
	err = storageManager.Migrate(dbVersion, migrationsDir)
	if err != nil {
		log.WithField("error", err).
			Fatal("failed to run migrations")
	}

	// setup API port
	a := api.New(storageManager)

	// setup http adapter
	h := server.NewAPIHandler(a)
	s := server.New(h)

	// start server
	log.WithFields(log.Fields{
		"address": conf.Addr,
		"port":    conf.Port,
	}).Info("server starting")

	log.WithField("error", s.Start(conf.Addr, conf.Port)).
		Fatal("server error")
}
