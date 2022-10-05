package config

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

var AppContext context.Context = nil

func onOsSignal(cleanupCb func()) {
	osSignal := make(chan os.Signal)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGTERM)

	select {
	case <- osSignal:
		log.Error().Msg("Caugth OS kill")
		cleanupCb()
		os.Exit(1)
	}
}


func init() {
	var cancel func()
	AppContext, cancel = context.WithCancel(context.Background())

	go onOsSignal(cancel)
}