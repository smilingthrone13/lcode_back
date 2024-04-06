package manager

import (
	"lcode/config"
	"lcode/internal/manager/problem_manager"
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

	Managers struct {
		ProblemManager *problem_manager.Manager
	}
)

func New(p *InitParams, services *service.Services) *Managers {
	problemManager := problem_manager.New(
		p.Config,
		p.Logger,
		p.TransactionManager,
		&problem_manager.Services{
			TaskService:         services.Task,
			TaskTemplateService: services.TaskTemplate,
			TestCaseService:     services.TestCase,
		},
	)

	return &Managers{
		ProblemManager: problemManager,
	}
}
