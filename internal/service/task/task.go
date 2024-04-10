package task

import (
	"context"
	"github.com/pkg/errors"
	"lcode/internal/domain"
	"log/slog"
)

type Service struct {
	logger     *slog.Logger
	repository TaskRepo
}

func New(
	logger *slog.Logger,
	repository TaskRepo,
) *Service {
	return &Service{logger: logger, repository: repository}
}

func (s *Service) Create(ctx context.Context, dto domain.TaskCreateInput) (taskID string, err error) {
	taskID, err = s.repository.Create(ctx, dto)
	if err != nil {
		return taskID, errors.Wrap(err, "Create Task service:")
	}

	return taskID, nil
}

func (s *Service) Update(ctx context.Context, id string, dto domain.TaskUpdateInput) error {
	err := s.repository.Update(ctx, id, dto)
	if err != nil {
		return errors.Wrap(err, "Update Task service:")
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, "Delete Task service:")
	}

	return nil
}

func (s *Service) GetByID(ctx context.Context, id string) (domain.Task, error) {
	t, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return domain.Task{}, errors.Wrap(err, "GetByID Task service:")
	}

	return t, nil
}

func (s *Service) GetAllByParams(ctx context.Context, params domain.TaskParams) (domain.TaskList, error) {
	tList, err := s.repository.GetAllByParams(ctx, params)
	if err != nil {
		return domain.TaskList{}, errors.Wrap(err, "GetAllByParams Task service:")
	}

	return tList, nil
}

func (s *Service) GetAvailableAttributes(ctx context.Context) (domain.TaskAttributes, error) {
	ta, err := s.repository.GetAvailableAttributes(ctx)
	if err != nil {
		return domain.TaskAttributes{}, errors.Wrap(err, "GetAvailableAttributes Task service:")
	}

	return ta, nil
}
