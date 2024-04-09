package comment

import (
	"context"
	"lcode/internal/domain"
)

type Comment interface {
	Create(ctx context.Context, dto domain.CommentCreateDTO) (domain.Comment, error)
	Update(ctx context.Context, dto domain.CommentUpdateDTO) (domain.Comment, error)
	Delete(ctx context.Context, id string) error

	GetThreadsByParamsAndEntityID(
		ctx context.Context,
		id string,
		params domain.CommentParams,
	) (domain.ThreadList, error)
}

type CommentRepo interface {
	Create(ctx context.Context, dto domain.CommentCreateDTO) (domain.Comment, error)
	Update(ctx context.Context, dto domain.CommentUpdateDTO) (domain.Comment, error)
	Delete(ctx context.Context, id string) error

	GetThreadsByParamsAndEntityID(
		ctx context.Context,
		id string,
		params domain.CommentParams,
	) (domain.ThreadList, error)
}
