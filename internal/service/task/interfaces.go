package task

import (
	"context"
	"lcode/internal/domain"
)

type Task interface {
	Create(ctx context.Context, dto domain.TaskCreateInput) (string, error)
	Update(ctx context.Context, id string, dto domain.TaskUpdateInput) error
	Delete(ctx context.Context, id string) error

	GetByID(ctx context.Context, id string) (domain.Task, error)
	GetAllByParams(ctx context.Context, params domain.TaskParams) (domain.TaskList, error)

	GetAvailableAttributes(ctx context.Context) (domain.TaskAttributes, error)
}

type TaskRepo interface {
	Create(ctx context.Context, dto domain.TaskCreateInput) (string, error)
	Update(ctx context.Context, id string, dto domain.TaskUpdateInput) error
	Delete(ctx context.Context, id string) error

	GetByID(ctx context.Context, id string) (domain.Task, error)
	GetAllByParams(ctx context.Context, params domain.TaskParams) (domain.TaskList, error)

	GetAvailableAttributes(ctx context.Context) (domain.TaskAttributes, error)
}
