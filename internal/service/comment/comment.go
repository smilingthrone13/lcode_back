package comment

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
	repository         CommentRepo
}

func New(
	logger *slog.Logger,
	transactionManager *postgres.TransactionProvider,
	repository CommentRepo,
) *Service {
	return &Service{logger: logger, transactionManager: transactionManager, repository: repository}
}

func (s *Service) Create(ctx context.Context, dto domain.CommentCreateDTO) (c domain.Comment, err error) {
	tx, err := s.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return c, errors.Wrap(err, "Comment Service Create:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	c, err = s.repository.Create(ctx, dto.Input)
	if err != nil {
		return c, errors.Wrap(err, "Comment Service Create:")
	}

	if err = tx.Commit(ctx); err != nil {
		return c, errors.Wrap(err, "Comment Service Create:")
	}

	return c, nil
}

func (s *Service) Update(ctx context.Context, dto domain.CommentUpdateDTO) (c domain.Comment, err error) {
	tx, err := s.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return c, errors.Wrap(err, "Comment Service Update:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	c, err = s.repository.Update(ctx, dto)
	if err != nil {
		return c, errors.Wrap(err, "Comment Service Update:")
	}

	if err = tx.Commit(ctx); err != nil {
		return c, errors.Wrap(err, "Comment Service Update:")
	}

	return c, nil
}

func (s *Service) Delete(ctx context.Context, dto domain.CommentDeleteDTO) error {
	tx, err := s.transactionManager.NewTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "Comment Service Delete:")
	}
	ctx = context.WithValue(ctx, postgres.TxKey{}, tx)
	defer tx.Rollback(ctx)

	err = s.repository.Delete(ctx, dto)
	if err != nil {
		return errors.Wrap(err, "Comment Service Delete:")
	}

	if err = tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "Comment Service Delete:")
	}

	return nil
}

func (s *Service) GetThreadsByParamsAndEntityID(
	ctx context.Context,
	dto domain.CommentParamsDTO,
) (tl domain.ThreadList, err error) {
	tl, err = s.repository.GetThreadsByParamsAndEntityID(ctx, dto.EntityID, dto.Input)
	if err != nil {
		return tl, errors.Wrap(err, "Comment Service GetThreadsByParamsAndEntityID:")
	}

	return tl, nil
}
