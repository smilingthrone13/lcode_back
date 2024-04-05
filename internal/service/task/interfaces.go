package task

import (
	"context"
	"lcode/internal/domain"
)

type Task interface {
	Create(ctx context.Context, dto domain.TaskCreate) (domain.Task, error)
	Update(ctx context.Context, id string, dto domain.TaskUpdate) (domain.Task, error)
	Delete(ctx context.Context, id string) (domain.Task, error)

	GetByID(ctx context.Context, id string) (domain.Task, error)
	GetAllByParams(ctx context.Context, params domain.TaskParams) (domain.TaskList, error)
}

type TaskRepo interface {
	Create(ctx context.Context, dto domain.TaskCreate) (domain.Task, error)
	Update(ctx context.Context, id string, dto domain.TaskUpdate) (domain.Task, error)
	Delete(ctx context.Context, id string) (domain.Task, error)

	GetByID(ctx context.Context, id string) (domain.Task, error)
	GetAllByParams(ctx context.Context, params domain.TaskParams) (domain.TaskList, error)
}
