package article

import (
	"context"
	"github.com/pkg/errors"
	"lcode/internal/domain"
	"log/slog"
)

type Service struct {
	logger     *slog.Logger
	repository ArticleRepo
}

func New(
	logger *slog.Logger,
	repository ArticleRepo,
) *Service {
	return &Service{logger: logger, repository: repository}
}

func (s *Service) Create(ctx context.Context, dto domain.ArticleCreateInput) (a domain.Article, err error) {
	a, err = s.repository.Create(ctx, dto)
	if err != nil {
		return a, errors.Wrap(err, "Article Service Create:")
	}

	return a, nil
}

func (s *Service) Update(ctx context.Context, dto domain.ArticleUpdateInput) (a domain.Article, err error) {
	a, err = s.repository.Update(ctx, dto)
	if err != nil {
		return a, errors.Wrap(err, "Article Service Update:")
	}

	return a, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, "Article Service Delete:")
	}

	return nil
}

func (s *Service) GetByID(ctx context.Context, id string) (a domain.Article, err error) {
	a, err = s.repository.GetByID(ctx, id)
	if err != nil {
		return a, errors.Wrap(err, "Article Service GetByID:")
	}

	return a, nil
}

func (s *Service) GetPracticeArticle(ctx context.Context) (a domain.Article, err error) {
	a, err = s.repository.GetPracticeArticle(ctx)
	if err != nil {
		return a, errors.Wrap(err, "Article Service GetPracticeArticle:")
	}

	return a, nil
}

func (s *Service) GetAllByParams(ctx context.Context, params domain.ArticleParams) (al domain.ArticleList, err error) {
	al, err = s.repository.GetAllByParams(ctx, params)
	if err != nil {
		return al, errors.Wrap(err, "Article Service GetAllByParams:")
	}

	return al, nil
}

func (s *Service) GetAvailableAttributes(ctx context.Context) (domain.ArticleAttributes, error) {
	aa, err := s.repository.GetAvailableAttributes(ctx)
	if err != nil {
		return aa, errors.Wrap(err, "Article Service GetAvailableAttributes:")
	}

	return aa, nil
}
