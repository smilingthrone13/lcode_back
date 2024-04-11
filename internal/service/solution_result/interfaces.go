package solution_result

import (
	"context"
	"lcode/internal/domain"
)

type (
	SolutionResult interface {
		CreateBatch(ctx context.Context, results ...domain.SolutionResult) error
		ResultsBySolutionID(ctx context.Context, solutionID string) ([]domain.SolutionResult, error)
	}

	SolutionResultRepo interface {
		CreateBatch(ctx context.Context, results ...domain.SolutionResult) error
		ResultsBySolutionID(ctx context.Context, solutionID string) ([]domain.SolutionResult, error)
	}
)
