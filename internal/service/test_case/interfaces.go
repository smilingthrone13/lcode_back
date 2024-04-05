package test_case

import (
	"context"
	"lcode/internal/domain"
)

type TestCase interface {
	Create(ctx context.Context, dto domain.TestCaseCreate) (domain.TestCase, error)
	Update(ctx context.Context, id string, dto domain.TestCaseUpdate) (domain.TestCase, error)
	Delete(ctx context.Context, id string) (domain.TestCase, error)

	GetByID(ctx context.Context, id string) (domain.TestCase, error)
	GetAllByTaskID(ctx context.Context, id string) ([]domain.TestCase, error)
}
