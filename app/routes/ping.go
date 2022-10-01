package routes

import (
	"net/http"
	"tinyurl/app/webserver"

	"github.com/rs/zerolog/log"
)

func Ping(srv webserver.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		if _, err := w.Write([]byte("Pong")); err != nil {
			log.Fatal().Msg("Failed to send pong!")
		}
	}
}