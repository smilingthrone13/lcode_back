package middleware

import (
	"lcode/config"
	"lcode/internal/handler/middleware/access"
	"lcode/internal/handler/middleware/auth"
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
		Access  *access.Middleware
		Auth    *auth.Middleware
		General *general.Middleware
	}
)

func New(p *InitParams, services *service.Services) *Middlewares {
	generalMiddleware := general.New(p.Config, p.Logger, &general.Services{})
	accessMiddleware := access.New(p.Config, p.Logger, &access.Services{
		Auth: services.Auth,
	})

	authMiddleware := auth.New(p.Config, p.Logger, &auth.Services{})

	return &Middlewares{
		Access:  accessMiddleware,
		Auth:    authMiddleware,
		General: generalMiddleware,
	}
}
