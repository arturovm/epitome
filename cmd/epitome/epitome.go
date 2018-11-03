package main

import (
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/arturovm/epitome/conf"
	"github.com/arturovm/epitome/data"
	"github.com/arturovm/epitome/router"
	log "github.com/sirupsen/logrus"
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
	err := data.Setup()
	if err != nil {
		log.WithField("error", err).
			Fatal("failed to initialize data directory")
	}

	// setup router
	r := router.Get()

	// start server
	log.WithFields(log.Fields{
		"address": conf.Addr,
		"port":    conf.Port,
	}).Info("server starting")

	addr := net.JoinHostPort(conf.Addr, strconv.Itoa(conf.Port))
	log.WithField("error", http.ListenAndServe(addr, r)).
		Fatal("server error")
}
