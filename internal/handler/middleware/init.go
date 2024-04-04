package middleware

import (
	"lcode/config"
	"lcode/internal/handler/middleware/general"
	"lcode/internal/service"
	"log/slog"
)

type (
	InitParams struct {
		Config *config.Config
		Logger *slog.Logger
	}

	Middlewares struct {
		General *general.Middleware
	}
)

func New(p *InitParams, services *service.Services) *Middlewares {
	generalMiddleware := general.New(p.Config, p.Logger, &general.Services{})

	return &Middlewares{
		General: generalMiddleware,
	}
}
