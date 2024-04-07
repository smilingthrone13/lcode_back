package task_template

import (
	"context"
	"lcode/internal/domain"
)

type TaskTemplate interface {
	Create(ctx context.Context, taskID string, dto domain.TaskTemplateCreateInput) (domain.TaskTemplate, error)
	Update(ctx context.Context, id string, dto domain.TaskTemplateUpdateInput) (domain.TaskTemplate, error)
	Delete(ctx context.Context, id string) error

	GetByID(ctx context.Context, id string) (domain.TaskTemplate, error)
	GetAllByTaskID(ctx context.Context, id string) ([]domain.TaskTemplate, error)
}

type TaskTemplateRepo interface {
	Create(ctx context.Context, taskID string, dto domain.TaskTemplateCreateInput) (domain.TaskTemplate, error)
	Update(ctx context.Context, id string, dto domain.TaskTemplateUpdateInput) (domain.TaskTemplate, error)
	Delete(ctx context.Context, id string) error

	GetByID(ctx context.Context, id string) (domain.TaskTemplate, error)
	GetAllByTaskID(ctx context.Context, id string) ([]domain.TaskTemplate, error)
}
