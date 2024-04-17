package article

import (
	"context"
	"lcode/internal/domain"
)

type Article interface {
	Create(ctx context.Context, dto domain.ArticleCreateInput) (domain.Article, error)
	Update(ctx context.Context, dto domain.ArticleUpdateInput) (domain.Article, error)
	Delete(ctx context.Context, id string) error

	CreateDefault(ctx context.Context, user domain.User) error

	GetByID(ctx context.Context, id string) (domain.Article, error)
	GetAllByParams(ctx context.Context, params domain.ArticleParams) (domain.ArticleList, error)

	GetAvailableAttributes(ctx context.Context) (domain.ArticleAttributes, error)
}

type ArticleRepo interface {
	Create(ctx context.Context, dto domain.ArticleCreateInput) (domain.Article, error)
	Update(ctx context.Context, dto domain.ArticleUpdateInput) (domain.Article, error)
	Delete(ctx context.Context, id string) error

	CreateDefault(ctx context.Context, user domain.User) error

	GetByID(ctx context.Context, id string) (domain.Article, error)
	GetAllByParams(ctx context.Context, params domain.ArticleParams) (domain.ArticleList, error)

	GetAvailableAttributes(ctx context.Context) (domain.ArticleAttributes, error)
}
