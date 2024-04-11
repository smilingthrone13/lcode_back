package solution

import (
	"context"
	"lcode/internal/domain"
)

type (
	Solution interface {
		Create(ctx context.Context, entity domain.CreateSolutionEntity) (sol domain.Solution, err error)
		Update(ctx context.Context, entity domain.UpdateSolutionDTO) (sol domain.Solution, err error)
		SolutionsByUserAndTask(ctx context.Context, dto domain.GetSolutionsDTO) ([]domain.Solution, error)
		SolutionByID(ctx context.Context, id string) (sol domain.Solution, err error)
	}

	SolutionRepo interface {
		Create(ctx context.Context, entity domain.CreateSolutionEntity) (sol domain.Solution, err error)
		Update(ctx context.Context, entity domain.UpdateSolutionDTO) (sol domain.Solution, err error)
		SolutionsByUserAndTask(ctx context.Context, dto domain.GetSolutionsDTO) ([]domain.Solution, error)
		SolutionByID(ctx context.Context, id string) (sol domain.Solution, err error)
	}
)
