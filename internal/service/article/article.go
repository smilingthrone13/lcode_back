package article

import (
	"context"
	"github.com/pkg/errors"
	"lcode/internal/domain"
	"lcode/pkg/postgres"
	"log/slog"
)

type Service struct {
	logger             *slog.Logger
	transactionManager *postgres.TransactionProvider
	repository         ArticleRepo
}

func New(
	logger *slog.Logger,
	transactionManager *postgres.TransactionProvider,
	repository ArticleRepo,
) *Service {
	return &Service{logger: logger, transactionManager: transactionManager, repository: repository}
}

func (s *Service) Create(ctx context.Context, dto domain.ArticleCreateInput) (a domain.Article, err error) {
	tx, err := s.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return a, errors.Wrap(err, "Article Service Create:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	a, err = s.repository.Create(ctx, dto)
	if err != nil {
		return a, errors.Wrap(err, "Article Service Create:")
	}

	if err = tx.Commit(ctx); err != nil {
		return a, errors.Wrap(err, "Article Service Create:")
	}

	return a, nil
}

func (s *Service) Update(ctx context.Context, dto domain.ArticleUpdateInput) (a domain.Article, err error) {
	tx, err := s.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return a, errors.Wrap(err, "Article Service Update:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	a, err = s.repository.Update(ctx, dto)
	if err != nil {
		return a, errors.Wrap(err, "Article Service Update:")
	}

	if err = tx.Commit(ctx); err != nil {
		return a, errors.Wrap(err, "Article Service Update:")
	}

	return a, nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	tx, err := s.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "Article Service Delete:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	err = s.repository.Delete(ctx, id)
	if err != nil {
		return errors.Wrap(err, "Article Service Delete:")
	}

	if err = tx.Commit(ctx); err != nil {
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
