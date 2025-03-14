package service

import (
	"lcode/config"
	"lcode/internal/infra/repository"
	"lcode/internal/service/article"
	"lcode/internal/service/auth"
	"lcode/internal/service/comment"
	"lcode/internal/service/solution"
	solutionResult "lcode/internal/service/solution_result"
	"lcode/internal/service/task"
	taskTemplate "lcode/internal/service/task_template"
	testCase "lcode/internal/service/test_case"
	"lcode/internal/service/thumbnails"
	"lcode/internal/service/user_fs"
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
		UserFS         user_fs.UserFS
		Thumbnails     thumbnails.Thumbnails
		Auth           auth.Authorization
		Task           task.Task
		TaskTemplate   taskTemplate.TaskTemplate
		TestCase       testCase.TestCase
		Solution       solution.Solution
		SolutionResult solutionResult.SolutionResult
		UserProgress   userProgress.UserProgress
		Article        article.Article
		Comment        comment.Comment
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
	articleService := article.New(p.Logger, p.TransactionManager, repos.Article)
	commentService := comment.New(p.Logger, p.TransactionManager, repos.Comment)
	thumbnailsService := thumbnails.New(p.Config, p.Logger)
	userFsService := user_fs.New(p.Config, p.Logger, &user_fs.Services{
		Thumbnails: thumbnailsService,
	})

	return &Services{
		Thumbnails:     thumbnailsService,
		UserFS:         userFsService,
		Auth:           authService,
		Task:           taskService,
		TaskTemplate:   taskTemplateService,
		TestCase:       testCaseService,
		Solution:       solutionService,
		SolutionResult: solutionResultService,
		UserProgress:   userProgressService,
		Article:        articleService,
		Comment:        commentService,
	}
}
