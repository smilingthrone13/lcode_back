package manager

import (
	"lcode/config"
	"lcode/internal/infra/webapi"
	"lcode/internal/manager/problem_manager"
	"lcode/internal/manager/solution_manager"
	"lcode/internal/manager/user_manager"
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
		UserManager     *user_manager.Manager
		ProblemManager  *problem_manager.Manager
		SolutionManager *solution_manager.Manager
	}
)

func New(p *InitParams, services *service.Services, apis *webapi.APIs) *Managers {
	problemManager := problem_manager.New(
		p.Config,
		p.Logger,
		p.TransactionManager,
		&problem_manager.Services{
			TaskService:         services.Task,
			TaskTemplateService: services.TaskTemplate,
			TestCaseService:     services.TestCase,
			Judge:               apis.Judge,
		},
	)

	solutionManager := solution_manager.New(
		p.Config,
		p.Logger,
		p.TransactionManager,
		&solution_manager.Services{
			ProblemManager: problemManager,
			Solution:       services.Solution,
			SolutionResult: services.SolutionResult,
			Judge:          apis.Judge,
		},
	)

	userManager := user_manager.New(
		p.Config,
		p.Logger,
		p.TransactionManager,
		&user_manager.Services{
			Auth:   services.Auth,
			UserFS: services.UserFS,
		},
	)

	return &Managers{
		UserManager:     userManager,
		ProblemManager:  problemManager,
		SolutionManager: solutionManager,
	}
}
