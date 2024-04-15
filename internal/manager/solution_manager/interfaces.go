package solution_manager

import (
	"context"
	"lcode/internal/domain"
)

type (
	SolutionManager interface {
		CreateSolution(ctx context.Context, dto domain.CreateSolutionDTO) (sol domain.Solution, err error)
		GetAvailableSolutionStatuses() ([]domain.JudgeStatusInfo, error)
	}

	ProblemManager interface {
		FullProblemByTaskID(ctx context.Context, taskID string) (domain.Problem, error)
	}

	Judge interface {
		CreateSubmission(
			ctx context.Context,
			data domain.CreateJudgeSubmission,
		) (domain.JudgeSubmissionInfo, error)

		GetAvailableStatuses(ctx context.Context) ([]domain.JudgeStatusInfo, error)
	}
)
