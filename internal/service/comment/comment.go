package comment

import (
	"context"
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

func (s Service) Create(ctx context.Context, dto domain.CommentCreateDTO) (domain.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Update(ctx context.Context, dto domain.CommentUpdateDTO) (domain.Comment, error) {
	//TODO implement me
	panic("implement me")
}

func (s Service) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) GetThreadsByParamsAndEntityID(
	ctx context.Context,
	id string,
	params domain.CommentParams,
) (domain.ThreadList, error) {
	//TODO implement me
	panic("implement me")
}
