package problem_manager

import (
	"golang.org/x/net/context"
	"lcode/internal/domain"
)

type ProblemManager interface {
	CreateProblem(ctx context.Context, dto domain.ProblemCreateInput) (domain.Problem, error)
	UpdateProblemTask(ctx context.Context, taskID string, dto domain.TaskUpdateInput) (domain.Problem, error)
	DeleteProblem(ctx context.Context, taskID string) error

	CreateProblemTaskTemplate(ctx context.Context, taskID string, dto domain.TaskTemplateCreateInput) (domain.Problem, error)
	UpdateProblemTaskTemplate(ctx context.Context, templateID string, dto domain.TaskTemplateUpdateInput) (domain.Problem, error)
	DeleteProblemTaskTemplate(ctx context.Context, templateID string) error

	CreateProblemTestCase(ctx context.Context, taskID string, dto domain.TestCaseCreateInput) (domain.Problem, error)
	UpdateProblemTestCase(ctx context.Context, caseID string, dto domain.TestCaseUpdateInput) (domain.Problem, error)
	DeleteProblemTestCase(ctx context.Context, caseID string) error

	FullProblemByTaskID(ctx context.Context, taskID string) (domain.Problem, error)
	TaskListByParams(ctx context.Context, dto domain.TaskParams) (domain.TaskList, error)

	GetAvailableTaskAttributes(ctx context.Context) (domain.TaskAttributes, error)
}

//type SubmissionManager interface {
//	SubmitSolution()
//	CheckSolutionStatus()
//	SolutionsByUserAndTaskID()
//}
