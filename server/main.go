package main

import (
	"tinyurl/server/config"


	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading configs")
	}

	appContext := config.AppContext

	server, err := NewServer(appContext, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating server")
	}

	server.Start()
	defer server.Close()
}
