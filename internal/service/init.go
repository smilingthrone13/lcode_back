package service

import (
	"lcode/config"
	"lcode/internal/infra/repository"
	"lcode/internal/service/authorization"
	"lcode/pkg/postgres"
	"log/slog"
)

type (
	InitParams struct {
		Config             *config.Config
		Logger             *slog.Logger
		TransactionManager *postgres.TransactionProvider
	}

	Services struct {
		Auth *authorization.Service
	}
)

func New(p *InitParams, repos *repository.Repositories) *Services {
	authService := authorization.NewService(p.Config)

	return &Services{
		Auth: authService,
	}
}
