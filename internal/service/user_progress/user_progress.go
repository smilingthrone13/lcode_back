package user_progress

import (
	"context"
	"github.com/pkg/errors"
	"lcode/internal/domain"
	"log/slog"
)

type Service struct {
	logger     *slog.Logger
	repository UserProgressRepo
}

func New(
	logger *slog.Logger,
	repository UserProgressRepo,
) *Service {
	return &Service{
		logger:     logger,
		repository: repository,
	}
}

func (s *Service) GetStatisticsByUserID(
	ctx context.Context,
	userID string,
	statType string,
) (us domain.UserStatistic, err error) {
	us, err = s.repository.StatisticsByUserID(ctx, userID, statType)
	if err != nil {
		return us, errors.Wrap(err, "User Progress Service StatisticsByUserID:")
	}

	return us, nil
}

func (s *Service) GetProgressByUserID(ctx context.Context, userID string) (up domain.UserProgress, err error) {
	up, err = s.repository.ProgressByUserID(ctx, userID)
	if err != nil {
		return up, errors.Wrap(err, "User Progress Service ProgressByUserID:")
	}

	return up, nil
}
