package main

import (
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/ArturoVM/epitome/conf"
	"github.com/ArturoVM/epitome/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	if conf.Help {
		conf.PrintHelp()
		os.Exit(0)
	}
	if conf.Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		log.Debug().Msg("debug mode enabled")
	}

	r := router.Get()

	log.Info().
		Str("address", conf.Addr).
		Int("port", conf.Port).
		Msg("server starting")
	log.Fatal().
		Err(http.ListenAndServe(net.JoinHostPort(conf.Addr, strconv.Itoa(conf.Port)), r)).
		Msg("server error")
}
