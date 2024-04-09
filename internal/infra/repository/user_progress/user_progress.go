package user_progress

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	sql_query_maker "github.com/m-a-r-a-t/sql-query-maker"
	"github.com/pkg/errors"
	"lcode/internal/domain"
	"lcode/pkg/postgres"
)

type Repository struct {
	db *postgres.DbManager
}

func New(db *postgres.DbManager) *Repository {
	return &Repository{db: db}
}

func (r *Repository) StatisticsByUserID(
	ctx context.Context,
	userID string,
	statType domain.StatisticsType,
) (s domain.UserStatistic, err error) {
	sq := sql_query_maker.NewQueryMaker(1)

	queryType := domain.StatisticDifficulty
	if statType != domain.StatisticDifficulty {
		queryType = domain.StatisticCategory
	}

	sq.Add(
		fmt.Sprintf(
			`
		WITH complete_s AS 
			(SELECT DISTINCT task_id, user_id, status
		                     FROM solution
		                     WHERE status = '%s'),
			statuses AS
			(SELECT DISTINCT s.task_id, s.user_id, COALESCE(complete_s.status, '%s') AS status
		      FROM solution s
		          LEFT JOIN complete_s
		              ON s.user_id = complete_s.user_id AND s.task_id = complete_s.task_id
		      WHERE s.user_id = ?)

		SELECT t.%s AS param, COUNT(s.task_id) AS count_done, COUNT(t.id) AS count_total
		FROM statuses s
    		RIGHT JOIN task t ON t.id = s.task_id
		GROUP BY param
		`,
			domain.ProgressCompleted, domain.ProgressInProgress, queryType),
		userID,
	)

	query, args := sq.Make()

	stats := []domain.StatisticData{}

	err = pgxscan.Select(ctx, r.db.TxOrDB(ctx), &stats, query, args...)
	if err != nil {
		return s, errors.Wrap(err, "StatisticsByUserID User Progress Repo:")
	}

	s = domain.UserStatistic{
		Type:       queryType,
		Statistics: stats,
	}

	return s, nil
}

func (r *Repository) ProgressByUserID(ctx context.Context, userID string) (p domain.UserProgress, err error) {
	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add(
		fmt.Sprintf(
			`
		WITH complete_s AS 
			(SELECT DISTINCT task_id, user_id, status
		                     FROM solution
		                     WHERE status = '%s'),
			statuses AS
			(SELECT DISTINCT s.task_id, s.user_id, COALESCE(complete_s.status, '%s') AS status
		      FROM solution s
		          LEFT JOIN complete_s
		              ON s.user_id = complete_s.user_id AND s.task_id = complete_s.task_id
		      WHERE s.user_id = ?)

		SELECT status, array_agg(task_id) as task_ids
		FROM statuses
		GROUP BY status
		`,
			domain.ProgressCompleted, domain.ProgressInProgress),
		userID,
	)

	query, args := sq.Make()

	statuses := []domain.ProgressData{}

	err = pgxscan.Select(ctx, r.db.TxOrDB(ctx), &statuses, query, args...)
	if err != nil {
		return p, errors.Wrap(err, "ProgressByUserID User Progress Repo:")
	}

	p.Progress = statuses

	return p, nil
}
