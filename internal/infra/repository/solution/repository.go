package solution

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	sql_query_maker "github.com/m-a-r-a-t/sql-query-maker"
	"github.com/pkg/errors"
	"lcode/internal/domain"
	"lcode/pkg/postgres"
)

func New(db *postgres.DbManager) *Repository {
	return &Repository{db: db}
}

type Repository struct {
	db *postgres.DbManager
}

func (r *Repository) Create(ctx context.Context, entity domain.CreateSolutionEntity) (sol domain.Solution, err error) {
	sq := sql_query_maker.NewQueryMaker(7)

	sq.Add(
		`INSERT INTO solution (user_id, code, status, task_id, language_id) 
			   VALUES (?, ?, ?, ?, ?) 
               RETURNING id, user_id, code, status, runtime, memory, task_id, language_id`,
		entity.UserID,
		entity.Code,
		entity.Status,
		entity.TaskID,
		entity.LanguageID,
	)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &sol, query, args...)
	if err != nil {
		return domain.Solution{}, errors.Wrap(err, "Create solution repo")
	}

	return sol, nil
}

func (r *Repository) Update(ctx context.Context, dto domain.UpdateSolutionDTO) (sol domain.Solution, err error) {
	sq := sql_query_maker.NewQueryMaker(3)

	sq.Add(`UPDATE solution SET`)

	if dto.Status != nil {
		sq.Add("status = ?,", *dto.Status)
	}

	sq.Where("id = ?", dto.ID)
	sq.Add("RETURNING id, user_id, code, status, runtime, memory, task_id, language_id")

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &sol, query, args...)
	if err != nil {
		return domain.Solution{}, errors.Wrap(err, "Update solution repo")
	}

	return sol, nil
}

func (r *Repository) GetSolutionsByUserIdAndTaskId(ctx context.Context, userID, taskID string) ([]domain.Solution, error) {
	sq := sql_query_maker.NewQueryMaker(1)

	results := []domain.Solution{}

	sq.Add(`
			SELECT id, user_id, code, status, runtime, memory, task_id, language_id
			FROM solution
			WHERE user_id = ? AND task_id = ?`,
		userID,
		taskID,
	)

	query, args := sq.Make()

	err := pgxscan.Select(ctx, r.db.TxOrDB(ctx), &results, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "GetSolutionsByUserIdAndTaskId solution repo")
	}

	return results, nil
}
