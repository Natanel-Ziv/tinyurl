package webserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"tinyurl/app/tinyurl"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type Server interface {
	tinyurl.TinyURL
}

type WebServer struct {
	tinyurl.TinyURL

	cfg       Config
	ctx       context.Context
	cancelCtx context.CancelFunc
	ls        net.ListenConfig
	router    *mux.Router
}

func New(cfg Config, parrentCtx context.Context) (*WebServer, error) {
	err := validateConfigs(cfg)
	if err != nil {
		return nil, err
	}

	ctx, cancelCtx := context.WithCancel(parrentCtx)

	return &WebServer{
		cfg:       cfg,
		ctx:       ctx,
		cancelCtx: cancelCtx,
	}, nil
}

func (webserver *WebServer) Start(binder func(s Server, r *mux.Router)) error {
	webserver.router = mux.NewRouter().StrictSlash(true)
	binder(webserver, webserver.router)

	listener, err := webserver.ls.Listen(webserver.ctx, "tcp", fmt.Sprintf(":%d", webserver.cfg.Port))
	if err != nil {
		return err
	}

	log.Info().Msgf("Listening on port: %d", webserver.cfg.Port)

	return http.Serve(listener, webserver.router)
}
