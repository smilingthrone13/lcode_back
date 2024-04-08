package service

import (
	"lcode/config"
	"lcode/internal/infra/repository"
	"lcode/internal/service/auth"
	"lcode/internal/service/solution"
	solutionResult "lcode/internal/service/solution_result"
	"lcode/internal/service/task"
	taskTemplate "lcode/internal/service/task_template"
	testCase "lcode/internal/service/test_case"
	userProgress "lcode/internal/service/user_progress"
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
		Auth           *auth.Service
		Task           *task.Service
		TaskTemplate   *taskTemplate.Service
		TestCase       *testCase.Service
		Solution       *solution.Service
		SolutionResult *solutionResult.Service
		UserProgress   *userProgress.Service
	}
)

func New(p *InitParams, repos *repository.Repositories) *Services {
	authService := auth.New(p.Config, repos.Auth)
	taskService := task.New(p.Logger, repos.Task)
	taskTemplateService := taskTemplate.New(p.Logger, repos.TaskTemplate)
	testCaseService := testCase.New(p.Logger, repos.TestCase)
	solutionResultService := solutionResult.New(p.Config, repos.SolutionResult)
	solutionService := solution.New(p.Config, repos.Solution)
	userProgressService := userProgress.New(p.Logger, repos.UserProgress)

	return &Services{
		Auth:           authService,
		Task:           taskService,
		TaskTemplate:   taskTemplateService,
		TestCase:       testCaseService,
		Solution:       solutionService,
		SolutionResult: solutionResultService,
		UserProgress:   userProgressService,
	}
}
