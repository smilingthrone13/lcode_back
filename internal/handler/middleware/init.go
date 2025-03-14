package middleware

import (
	"lcode/config"
	"lcode/internal/handler/middleware/access"
	"lcode/internal/handler/middleware/article"
	"lcode/internal/handler/middleware/auth"
	"lcode/internal/handler/middleware/comment"
	"lcode/internal/handler/middleware/problem"
	"lcode/internal/handler/middleware/solution"
	userProgress "lcode/internal/handler/middleware/user_progress"
	"lcode/internal/manager"
	"lcode/internal/service"
	"log/slog"
)

type (
	InitParams struct {
		Config *config.Config
		Logger *slog.Logger
	}

	Middlewares struct {
		Access       *access.Middleware
		Auth         *auth.Middleware
		Problem      *problem.Middleware
		UserProgress *userProgress.Middleware
		Article      *article.Middleware
		Solution     *solution.Middleware
		Comment      *comment.Middleware
	}
)

func New(p *InitParams, services *service.Services, managers *manager.Managers) *Middlewares {
	accessMiddleware := access.New(
		p.Config,
		p.Logger,
		&access.Services{
			Auth: services.Auth,
		},
	)

	authMiddleware := auth.New(
		p.Config,
		p.Logger,
		&auth.Services{},
	)

	problemMiddleware := problem.New(
		p.Config,
		p.Logger,
		&problem.Managers{
			Problem: managers.ProblemManager,
		},
	)

	userProgressMiddleware := userProgress.New(
		p.Config,
		p.Logger,
	)

	articleMiddleware := article.New(
		p.Config,
		p.Logger,
	)

	solutionMiddleware := solution.New(
		p.Config,
		p.Logger,
		&solution.Services{
			Solution: services.Solution,
		},
	)

	commentMiddleware := comment.New(
		p.Config,
		p.Logger,
	)

	return &Middlewares{
		Access:       accessMiddleware,
		Auth:         authMiddleware,
		Problem:      problemMiddleware,
		UserProgress: userProgressMiddleware,
		Article:      articleMiddleware,
		Solution:     solutionMiddleware,
		Comment:      commentMiddleware,
	}
}
