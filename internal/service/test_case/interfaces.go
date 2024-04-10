package test_case

import (
	"context"
	"lcode/internal/domain"
)

type TestCase interface {
	Create(ctx context.Context, taskID string, dto domain.TestCaseCreateInput) error
	Update(ctx context.Context, id string, dto domain.TestCaseUpdateInput) error
	Delete(ctx context.Context, id string) error

	GetAllByTaskID(ctx context.Context, id string) ([]domain.TestCase, error)
}

type TestCaseRepo interface {
	Create(ctx context.Context, taskID string, dto domain.TestCaseCreateInput) error
	Update(ctx context.Context, id string, dto domain.TestCaseUpdateInput) error
	Delete(ctx context.Context, id string) error

	GetAllByTaskID(ctx context.Context, id string) ([]domain.TestCase, error)
}
