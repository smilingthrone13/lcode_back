package test_case

import (
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	sql_query_maker "github.com/m-a-r-a-t/sql-query-maker"
	"github.com/pkg/errors"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/pkg/postgres"
)

type Repository struct {
	cfg *config.Config
	db  *postgres.DbManager
}

func New(cfg *config.Config, db *postgres.DbManager) *Repository {
	return &Repository{cfg: cfg, db: db}
}

func (r *Repository) Create(ctx context.Context, taskID string, dto domain.TestCaseCreateInput) error {
	sq := sql_query_maker.NewQueryMaker(3)

	sq.Add(
		`
	INSERT INTO test_case (task_id, input, output)
	VALUES (?, ?, ?)
	`,
		taskID, dto.Input, dto.Output,
	)

	query, args := sq.Make()

	_, err := r.db.TxOrDB(ctx).Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "Create TestCase repo:")
	}

	return nil
}

func (r *Repository) Update(ctx context.Context, id string, dto domain.TestCaseUpdateInput) error {
	sq := sql_query_maker.NewQueryMaker(3)

	sq.Add("UPDATE test_case SET")

	if dto.Input != nil {
		sq.Add("input = ?", *dto.Input)
	}

	if dto.Output != nil {
		sq.Add("output = ?", *dto.Output)
	}

	sq.Add("WHERE id = ?", id)

	query, args := sq.Make()

	res, err := r.db.TxOrDB(ctx).Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "Update TestCase repo:")
	}

	if res.RowsAffected() == 0 {
		err = errors.New("TestCase not found!")

		return errors.Wrap(err, "Update TestCase repo:")
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add("DELETE FROM test_case WHERE id = ?", id)

	query, args := sq.Make()

	res, err := r.db.TxOrDB(ctx).Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "Delete TestCase repo:")
	}

	if res.RowsAffected() == 0 {
		err = errors.New("TestCase not found!")

		return errors.Wrap(err, "Delete TestCase repo:")
	}

	return nil
}

func (r *Repository) GetAllByTaskID(ctx context.Context, id string) ([]domain.TestCase, error) {
	tcs := []domain.TestCase{}

	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add("SELECT id, task_id, row_number() over (ORDER BY created_at) AS number, input, output FROM test_case WHERE task_id = ?", id)

	query, args := sq.Make()

	err := pgxscan.Select(ctx, r.db.TxOrDB(ctx), &tcs, query, args...)
	if err != nil {
		return tcs, errors.Wrap(err, "GetAllByTaskID TestCase repo:")
	}

	return tcs, nil
}
