package general

import (
	"lcode/config"
	"log/slog"
)

type (
	Services struct {
	}

	Middleware struct {
		cfg      *config.Config
		logger   *slog.Logger
		services *Services
	}
)

func New(cfg *config.Config, logger *slog.Logger, services *Services) *Middleware {
	return &Middleware{
		cfg:      cfg,
		logger:   logger,
		services: services,
	}
}
