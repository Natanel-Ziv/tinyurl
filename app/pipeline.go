package main

import (
	"net/http"
	"tinyurl/app/routes"
	"tinyurl/app/webserver"

	"github.com/gorilla/mux"
)

func BuildPipeline(srv webserver.Server, r *mux.Router) {
	r.HandleFunc("/ping", routes.Ping(srv)).Methods(http.MethodGet)
}