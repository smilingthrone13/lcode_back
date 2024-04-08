package solution_manager

import (
	"context"
	"lcode/internal/domain"
)

type SolutionManager interface {
	CreateSolution(ctx context.Context, dto domain.CreateSolutionDTO) (sol domain.Solution, err error)
}
