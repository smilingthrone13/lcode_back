package solution

import (
	"context"
	"lcode/internal/domain"
)

type (
	Solution interface {
		Create(ctx context.Context, entity domain.CreateSolutionDTO) (sol domain.Solution, err error)
		Update(ctx context.Context, entity domain.UpdateSolutionDTO) (sol domain.Solution, err error)
		GetSolutionsByUserIdAndTaskId(ctx context.Context, userID, taskID string) ([]domain.Solution, error)
	}

	SolutionRepo interface {
		Create(ctx context.Context, entity domain.CreateSolutionDTO) (sol domain.Solution, err error)
		Update(ctx context.Context, entity domain.UpdateSolutionDTO) (sol domain.Solution, err error)
		GetSolutionsByUserIdAndTaskId(ctx context.Context, userID, taskID string) ([]domain.Solution, error)
	}
)
