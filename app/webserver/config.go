package webserver

import "errors"

type Config struct {
	Port int
}

func validateConfigs(cfg Config) error {
	if cfg.Port < 0 {
		return errors.New("Port must be > 0")
	}
	return nil
}