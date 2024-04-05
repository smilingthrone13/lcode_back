package handler

import (
	"lcode/config"
	authH "lcode/internal/handler/http/auth"
	generalH "lcode/internal/handler/http/general"
	"lcode/internal/service"
	"lcode/pkg/postgres"
	"log/slog"
)

type (
	InitParams struct {
		Config             *config.Config
		Logger             *slog.Logger
		TransactionManager *postgres.TransactionProvider
	}

	HTTPHandlers struct {
		General *generalH.Handler
		Auth    *authH.Handler
	}

	Handlers struct {
		HTTP *HTTPHandlers
	}
)

func New(p *InitParams, services *service.Services) *Handlers {
	generalHandler := generalH.New(p.Config, p.Logger, &generalH.Services{})
	authHandler := authH.New(p.Config, p.Logger, &authH.Services{
		Auth: services.Auth,
	})

	return &Handlers{
		&HTTPHandlers{
			General: generalHandler,
			Auth:    authHandler,
		},
	}
}
