package user_progress

import (
	"context"
	"lcode/internal/domain"
)

type UserProgress interface {
	GetStatisticsByUserID(ctx context.Context, userID string, statType domain.StatisticsType) (domain.UserStatistic, error)
	GetProgressByUserID(ctx context.Context, userID string) (domain.UserProgress, error)
}

type UserProgressRepo interface {
	StatisticsByUserID(ctx context.Context, userID string, statType domain.StatisticsType) (domain.UserStatistic, error)
	ProgressByUserID(ctx context.Context, userID string) (domain.UserProgress, error)
}
