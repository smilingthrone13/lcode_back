package solution_result

import (
	"context"
	"github.com/pkg/errors"
	"lcode/config"
	"lcode/internal/domain"
)

type (
	Service struct {
		config     *config.Config
		repository SolutionResultRepo
	}
)

func New(conf *config.Config, repository SolutionResultRepo) *Service {
	return &Service{
		config:     conf,
		repository: repository,
	}
}

func (s *Service) CreateBatch(ctx context.Context, results ...domain.SolutionResult) error {
	err := s.repository.CreateBatch(ctx, results...)
	if err != nil {
		return errors.Wrap(err, "CreateBatch solution_result service")
	}

	return nil
}

func (s *Service) ResultsBySolutionID(ctx context.Context, solutionID string) ([]domain.SolutionResult, error) {
	results, err := s.repository.ResultsBySolutionID(ctx, solutionID)
	if err != nil {
		return nil, errors.Wrap(err, "ResultsBySolutionID solution_result service")
	}

	return results, nil
}
