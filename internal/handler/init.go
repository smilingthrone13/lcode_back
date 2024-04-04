package handler

import (
	"lcode/config"
	generalH "lcode/internal/handler/http/general"
	"lcode/internal/service"
	"log/slog"
)

type (
	InitParams struct {
		Config             *config.Config
		Logger             *slog.Logger
		TransactionManager *postgres.TransactionProvider
		ConnectionPool     *connection.Pool
	}

	HTTPHandlers struct {
		General *generalH.Handler
	}

	Handlers struct {
		HTTP *HTTPHandlers
	}
)

func New(p *InitParams, services *service.Services) *Handlers {
	generalHandler := generalH.New(p.Config, p.Logger, &generalH.Services{
		Auth: services.Auth,
	})

	return &Handlers{
		&HTTPHandlers{
			General: generalHandler,
		},
	}
}
