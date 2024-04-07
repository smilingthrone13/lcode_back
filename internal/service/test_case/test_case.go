package test_case

import (
	"context"
	"github.com/pkg/errors"
	"lcode/internal/domain"
	"log/slog"
)

type Service struct {
	logger     *slog.Logger
	repository TestCaseRepo
}

func New(
	logger *slog.Logger,
	repository TestCaseRepo,
) *Service {
	return &Service{logger: logger, repository: repository}
}

func (s *Service) Create(ctx context.Context, taskID string, dto domain.TestCaseCreateInput) (domain.TestCase, error) {
	tc, err := s.repository.Create(ctx, taskID, dto)
	if err != nil {
		return domain.TestCase{}, errors.Wrap(err, "Create TestCase service:")
	}

	return tc, nil
}

func (s *Service) Update(ctx context.Context, id string, dto domain.TestCaseUpdateInput) (domain.TestCase, error) {
	tc, err := s.repository.Update(ctx, id, dto)
	if err != nil {
		return domain.TestCase{}, errors.Wrap(err, "Update TestCase service:")
	}

	return tc, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, "Delete TestCase service:")
	}

	return nil
}

func (s *Service) GetByID(ctx context.Context, id string) (domain.TestCase, error) {
	tc, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.TestCase{}, errors.Wrap(err, "GetByID TestCase service:")
	}

	return tc, nil
}

func (s *Service) GetAllByTaskID(ctx context.Context, id string) ([]domain.TestCase, error) {
	tcs, err := s.repository.GetAllByTaskID(ctx, id)
	if err != nil {
		return []domain.TestCase{}, errors.Wrap(err, "GetAllByTaskID TestCase service:")
	}

	return tcs, nil
}
