package main

import (
	"errors"
	"net/http"
	"os"
	"tinyurl/app/utils"
	"tinyurl/app/webserver"

	"github.com/rs/zerolog/log"
)

func main() {
	appContext := utils.AppContext


	configs := webserver.Config{
		Port: 1338,
	}

	srv, err := webserver.New(configs, appContext)

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

