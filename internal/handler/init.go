package handler

import (
	"lcode/config"
	authH "lcode/internal/handler/http/auth"
	problemH "lcode/internal/handler/http/problem"
	userProgressH "lcode/internal/handler/http/user_progress"
	"lcode/internal/manager"
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
		Auth         *authH.Handler
		Problem      *problemH.Handler
		UserProgress *userProgressH.Handler
	}

	Handlers struct {
		HTTP *HTTPHandlers
	}
)

func New(p *InitParams, services *service.Services, managers *manager.Managers) *Handlers {
	authHandler := authH.New(
		p.Config,
		p.Logger,
		&authH.Services{
			Auth: services.Auth,
		},
	)

	problemHandler := problemH.New(
		p.Config,
		p.Logger,
		&problemH.Managers{
			Problem: managers.ProblemManager,
		},
	)

	userProgressHandler := userProgressH.New(
		p.Config,
		p.Logger,
		&userProgressH.Services{
			UserProgress: services.UserProgress,
		},
	)

	return &Handlers{
		&HTTPHandlers{
			Auth:         authHandler,
			Problem:      problemHandler,
			UserProgress: userProgressHandler,
		},
	}
}
