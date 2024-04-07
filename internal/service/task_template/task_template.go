package task_template

import (
	"context"
	"github.com/pkg/errors"
	"lcode/internal/domain"
	"log/slog"
)

type Service struct {
	logger     *slog.Logger
	repository TaskTemplateRepo
}

func New(
	logger *slog.Logger,
	repository TaskTemplateRepo,
) *Service {
	return &Service{logger: logger, repository: repository}
}

func (s Service) Create(ctx context.Context, taskID string, dto domain.TaskTemplateCreateInput) (domain.TaskTemplate, error) {
	tt, err := s.repository.Create(ctx, taskID, dto)
	if err != nil {
		return domain.TaskTemplate{}, errors.Wrap(err, "Create TaskTemplate service:")
	}

	return tt, nil
}

func (s Service) Update(ctx context.Context, id string, dto domain.TaskTemplateUpdateInput) (domain.TaskTemplate, error) {
	tt, err := s.repository.Update(ctx, id, dto)
	if err != nil {
		return domain.TaskTemplate{}, errors.Wrap(err, "Update TaskTemplate service:")
	}

	return tt, nil
}

func (s Service) Delete(ctx context.Context, id string) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, "Delete TaskTemplate service:")
	}

	return nil
}

func (s Service) GetByID(ctx context.Context, id string) (domain.TaskTemplate, error) {
	tt, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.TaskTemplate{}, errors.Wrap(err, "GetByID TaskTemplate service:")
	}

	return tt, nil
}

func (s Service) GetAllByTaskID(ctx context.Context, id string) ([]domain.TaskTemplate, error) {
	tts, err := s.repository.GetAllByTaskID(ctx, id)
	if err != nil {
		return []domain.TaskTemplate{}, errors.Wrap(err, "GetAllByTaskID TaskTemplate service:")
	}

	return tts, nil
}
