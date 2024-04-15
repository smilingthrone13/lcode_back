package problem_manager

import (
	"golang.org/x/net/context"
	"lcode/internal/domain"
)

type ProblemManager interface {
	CreateProblem(ctx context.Context, dto domain.ProblemCreateDTO) (domain.Problem, error)
	UpdateProblemTask(ctx context.Context, dto domain.TaskUpdateDTO) (domain.Problem, error)
	DeleteProblem(ctx context.Context, taskID string) error

	CreateProblemTaskTemplate(ctx context.Context, dto domain.TaskTemplateCreateDTO) (domain.Problem, error)
	UpdateProblemTaskTemplate(ctx context.Context, dto domain.TaskTemplateUpdateDTO) (domain.Problem, error)
	DeleteProblemTaskTemplate(ctx context.Context, templateID string) error

	CreateProblemTestCase(ctx context.Context, dto domain.TestCaseCreateDTO) (domain.Problem, error)
	UpdateProblemTestCase(ctx context.Context, dto domain.TestCaseUpdateDTO) (domain.Problem, error)
	DeleteProblemTestCase(ctx context.Context, caseID string) error

	FullProblemByTaskID(ctx context.Context, taskID string) (domain.Problem, error)
	TaskListByParams(ctx context.Context, dto domain.TaskParams) (domain.TaskList, error)

	GetAvailableTaskAttributes(ctx context.Context) (domain.TaskAttributes, error)
	GetAvailableTaskLanguages() ([]domain.JudgeLanguageInfo, error)
}

type Judge interface {
	GetAvailableLanguages(ctx context.Context) ([]domain.JudgeLanguageInfo, error)
}
