package task_template

import (
	"context"
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

func (s Service) Create(ctx context.Context, dto domain.TaskTemplateCreate) (domain.TaskTemplate, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Update(ctx context.Context, id string, dto domain.TaskTemplateUpdate) (domain.TaskTemplate, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Delete(ctx context.Context, id string) (domain.TaskTemplate, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetByID(ctx context.Context, id string) (domain.TaskTemplate, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetAllByTaskID(ctx context.Context, id string) ([]domain.TaskTemplate, error) {
	//TODO implement me
	panic("implement me")
}
