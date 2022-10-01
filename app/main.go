package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"tinyurl/app/webserver"

	"github.com/rs/zerolog/log"
)

func main() {
	configs := webserver.Config{
		Port: 1338,
	}

	ctx := context.Background()

	srv, err := webserver.New(configs, ctx)

	if err != nil {
		log.Fatal().Msgf("Failed to start server: %+v", err)
		os.Exit(1)
	}

	err = srv.Start(BuildPipeline)
	if errors.Is(err, http.ErrServerClosed) {
		log.Warn().Msg("Server is shuting down")
	} else {
		log.Fatal().Msgf("Failure in server: %+v", err)
		os.Exit(1)
	}
}
