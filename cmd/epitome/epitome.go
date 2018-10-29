package main

import (
	"net/http"
	"os"

	"github.com/ArturoVM/epitome/router"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	log.Debug().Msg("debug mode enabled")

	r := router.Get()

	log.Info().
		Int("port", 8080).
		Msg("server starting")
	log.Fatal().
		Err(http.ListenAndServe(":8080", r)).
		Msg("server error")
}
