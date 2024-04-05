package task_template

import (
	"context"
	"lcode/internal/domain"
)

type TaskTemplate interface {
	Create(ctx context.Context, dto domain.TaskTemplateCreate) (domain.TaskTemplate, error)
	Update(ctx context.Context, id string, dto domain.TaskTemplateUpdate) (domain.TaskTemplate, error)
	Delete(ctx context.Context, id string) (domain.TaskTemplate, error)

	GetByID(ctx context.Context, id string) (domain.TaskTemplate, error)
	GetAllByTaskID(ctx context.Context, id string) ([]domain.TaskTemplate, error)
}

type TaskTemplateRepo interface {
	Create(ctx context.Context, dto domain.TaskTemplateCreate) (domain.TaskTemplate, error)
	Update(ctx context.Context, id string, dto domain.TaskTemplateUpdate) (domain.TaskTemplate, error)
	Delete(ctx context.Context, id string) (domain.TaskTemplate, error)

	GetByID(ctx context.Context, id string) (domain.TaskTemplate, error)
	GetAllByTaskID(ctx context.Context, id string) ([]domain.TaskTemplate, error)
}
