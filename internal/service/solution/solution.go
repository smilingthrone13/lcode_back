package solution

import (
	"context"
	"github.com/pkg/errors"
	"lcode/config"
	"lcode/internal/domain"
)

type (
	Service struct {
		config     *config.Config
		repository SolutionRepo
	}
)

func New(conf *config.Config, repository SolutionRepo) *Service {
	return &Service{
		config:     conf,
		repository: repository,
	}
}

func (s *Service) Create(ctx context.Context, entity domain.CreateSolutionEntity) (sol domain.Solution, err error) {
	sol, err = s.repository.Create(ctx, entity)
	if err != nil {
		return sol, errors.Wrap(err, "Create solution service")
	}

	return sol, nil
}

func (s *Service) Update(ctx context.Context, entity domain.UpdateSolutionDTO) (sol domain.Solution, err error) {
	sol, err = s.repository.Update(ctx, entity)
	if err != nil {
		return sol, errors.Wrap(err, "Update solution service")
	}

	return sol, nil
}

func (s *Service) SolutionsByUserAndTask(ctx context.Context, dto domain.GetSolutionsDTO) ([]domain.Solution, error) {
	solutions, err := s.repository.SolutionsByUserAndTask(ctx, dto)
	if err != nil {
		return nil, errors.Wrap(err, "SolutionsByUserAndTask solution service")
	}

	return solutions, nil
}

func (s *Service) SolutionByID(ctx context.Context, id string) (sol domain.Solution, err error) {
	sol, err = s.repository.SolutionByID(ctx, id)
	if err != nil {
		return domain.Solution{}, errors.Wrap(err, "SolutionByID solution service")
	}

	return sol, nil
}
