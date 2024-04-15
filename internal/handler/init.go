package handler

import (
	"lcode/config"
	articleH "lcode/internal/handler/http/article"
	authH "lcode/internal/handler/http/auth"
	commentH "lcode/internal/handler/http/comment"
	problemH "lcode/internal/handler/http/problem"
	solutionH "lcode/internal/handler/http/solution"
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
		Article      *articleH.Handler
		Solution     *solutionH.Handler
		Comment      *commentH.Handler
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
			Auth:        services.Auth,
			UserManager: managers.UserManager,
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

	articleHandler := articleH.New(
		p.Config,
		p.Logger,
		&articleH.Services{
			Article: services.Article,
		},
	)

	solutionHandler := solutionH.New(
		p.Config,
		p.Logger,
		&solutionH.Services{
			SolutionManager:       managers.SolutionManager,
			SolutionService:       services.Solution,
			SolutionResultService: services.SolutionResult,
		},
	)

	commentHandler := commentH.New(
		p.Config,
		p.Logger,
		&commentH.Services{
			Comment: services.Comment,
		},
	)

	return &Handlers{
		&HTTPHandlers{
			Auth:         authHandler,
			Problem:      problemHandler,
			UserProgress: userProgressHandler,
			Article:      articleHandler,
			Solution:     solutionHandler,
			Comment:      commentHandler,
		},
	}
}
