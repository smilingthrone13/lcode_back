package repository

import (
	"lcode/config"
	"lcode/internal/infra/repository/auth"
	"lcode/internal/infra/repository/solution"
	solutionResult "lcode/internal/infra/repository/solution_result"
	"lcode/internal/infra/repository/task"
	taskTemplate "lcode/internal/infra/repository/task_template"
	testCase "lcode/internal/infra/repository/test_case"
	userProgress "lcode/internal/infra/repository/user_progress"
	"lcode/pkg/postgres"
)

type (
	InitParams struct {
		Config *config.Config
		DB     *postgres.DbManager
	}

	Repositories struct {
		Auth           *auth.Repository
		Task           *task.Repository
		TaskTemplate   *taskTemplate.Repository
		TestCase       *testCase.Repository
		Solution       *solution.Repository
		SolutionResult *solutionResult.Repository
		UserProgress   *userProgress.Repository
	}
)

func New(p *InitParams) *Repositories {
	return &Repositories{
		Auth:           auth.New(p.DB),
		Task:           task.New(p.Config, p.DB),
		TaskTemplate:   taskTemplate.New(p.Config, p.DB),
		TestCase:       testCase.New(p.Config, p.DB),
		Solution:       solution.New(p.DB),
		SolutionResult: solutionResult.New(p.DB),
		UserProgress:   userProgress.New(p.DB),
	}
}
