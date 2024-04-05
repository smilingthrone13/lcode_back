package repository

import (
	"lcode/config"
	"lcode/internal/infra/repository/auth"
	"lcode/internal/infra/repository/task"
	taskTemplate "lcode/internal/infra/repository/task_template"
	testCase "lcode/internal/infra/repository/test_case"
	"lcode/pkg/postgres"
)

type (
	InitParams struct {
		Config *config.Config
		DB     *postgres.DbManager
	}

	Repositories struct {
		Auth         *auth.Repository
		Task         *task.Repository
		TaskTemplate *taskTemplate.Repository
		TestCase     *testCase.Repository
	}
)

func New(p *InitParams) *Repositories {
	return &Repositories{
		Auth:         auth.New(p.DB),
		Task:         task.New(p.Config, p.DB),
		TaskTemplate: taskTemplate.New(p.Config, p.DB),
		TestCase:     testCase.New(p.Config, p.DB),
	}
}
