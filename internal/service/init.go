package service

import (
	"lcode/config"
	"lcode/internal/infra/repository"
	"lcode/internal/service/authorization"
	"lcode/internal/service/task"
	taskTemplate "lcode/internal/service/task_template"
	testCase "lcode/internal/service/test_case"
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
		Auth         *authorization.Service
		Task         *task.Service
		TaskTemplate *taskTemplate.Service
		TestCase     *testCase.Service
	}
)

func New(p *InitParams, repos *repository.Repositories) *Services {
	authService := authorization.New(p.Config, repos.Auth)
	taskService := task.New(p.Logger, repos.Task)
	taskTemplateService := taskTemplate.New(p.Logger, repos.TaskTemplate)
	testCaseService := testCase.New(p.Logger, repos.TestCase)

	return &Services{
		Auth:         authService,
		Task:         taskService,
		TaskTemplate: taskTemplateService,
		TestCase:     testCaseService,
	}
}
