package article

import (
	"context"
	"lcode/internal/domain"
)

type Article interface {
	Create(ctx context.Context, dto domain.ArticleCreateInput) (domain.Article, error)
	Update(ctx context.Context, dto domain.ArticleUpdateInput) (domain.Article, error)
	Delete(ctx context.Context, id string) error

	GetByID(ctx context.Context, id string) (domain.Article, error)
	GetPracticeArticle(ctx context.Context) (domain.Article, error)
	GetAllByParams(ctx context.Context, params domain.ArticleParams) (domain.ArticleList, error)
}

type ArticleRepo interface {
	Create(ctx context.Context, dto domain.ArticleCreateInput) (domain.Article, error)
	Update(ctx context.Context, dto domain.ArticleUpdateInput) (domain.Article, error)
	Delete(ctx context.Context, id string) error

	GetByID(ctx context.Context, id string) (domain.Article, error)
	GetPracticeArticle(ctx context.Context) (domain.Article, error)
	GetAllByParams(ctx context.Context, params domain.ArticleParams) (domain.ArticleList, error)
}
