package task

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgconn"
	sql_query_maker "github.com/m-a-r-a-t/sql-query-maker"
	"github.com/pkg/errors"
	"lcode/config"
	"lcode/internal/domain"
	"lcode/pkg/db"
	"lcode/pkg/postgres"
	"lcode/pkg/struct_errors"
)

type Repository struct {
	cfg *config.Config
	db  *postgres.DbManager
}

func New(cfg *config.Config, db *postgres.DbManager) *Repository {
	return &Repository{cfg: cfg, db: db}
}

func (r *Repository) Create(ctx context.Context, dto domain.TaskCreateInput) (t domain.Task, err error) {
	sq := sql_query_maker.NewQueryMaker(6)

	sq.Add(
		`
	INSERT INTO task (name, description, difficulty, category, runtime_limit, memory_limit)
	VALUES (?, ?, ?, ?, ?, ?)
	RETURNING id, number, name, description, difficulty, category, runtime_limit, memory_limit
	`,
		dto.Name, dto.Description, dto.Difficulty, dto.Category, dto.RuntimeLimit, dto.MemoryLimit,
	)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &t, query, args...)
	if err == nil {
		return t, nil
	}

	var pgError *pgconn.PgError
	if ok := errors.As(err, &pgError); !ok {
		return t, errors.Wrap(err, "Create Task repo:")
	}

	switch pgError.Code {
	case postgres.ERRCODE_UNIQUE_VIOLATION:
		err = &struct_errors.ErrExist{Err: err, Msg: "Name already exist"}
	}

	return t, errors.Wrap(err, "Create Task repo:")

}

func (r *Repository) Update(ctx context.Context, id string, dto domain.TaskUpdateInput) (t domain.Task, err error) {
	sq := sql_query_maker.NewQueryMaker(7)

	sq.Add("UPDATE task SET")

	if dto.Name != nil {
		sq.Add("name = ?,", *dto.Name)
	}

	if dto.Description != nil {
		sq.Add("description = ?,", *dto.Description)
	}

	if dto.Category != nil {
		sq.Add("category = ?,", *dto.Category)
	}

	if dto.Difficulty != nil {
		sq.Add("difficulty = ?,", *dto.Difficulty)
	}

	if dto.RuntimeLimit != nil {
		sq.Add("runtime_limit = ?,", *dto.RuntimeLimit)
	}

	if dto.MemoryLimit != nil {
		sq.Add("memory_limit = ?,", *dto.MemoryLimit)
	}

	sq.Where("id = ?", id)
	sq.Add("RETURNING id, number, name, description, category, difficulty, runtime_limit, memory_limit")

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &t, query, args...)
	if err != nil {
		return t, errors.Wrap(err, "Update Task repo:")
	}

	return t, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add("DELETE FROM task WHERE id = ?", id)

	query, args := sq.Make()

	res, err := r.db.TxOrDB(ctx).Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "Delete Task repo:")
	}

	if res.RowsAffected() == 0 {
		err = errors.New("Task not found!")

		return errors.Wrap(err, "Delete Task repo:")
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (t domain.Task, err error) {
	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add(
		`
	SELECT id, number, name, description, category, difficulty, runtime_limit, memory_limit
	FROM task
	WHERE id = ?
	`,
		id)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &t, query, args...)
	if err != nil {
		return t, errors.Wrap(err, "GetByID Task repo:")
	}

	return t, nil
}

func (r *Repository) GetAllByParams(ctx context.Context, params domain.TaskParams) (tList domain.TaskList, err error) {
	tasks := []domain.Task{}
	sq := newFilter(r.cfg, 15)

	sq.Add("SELECT id, number, name, description, category, difficulty, runtime_limit, memory_limit FROM task t")

	if params.Pagination.AfterID != nil {
		q := fmt.Sprintf(
			"WHERE number %s (SELECT number FROM task WHERE id = ?)",
			db.GetLetterGreaterOrLessBySortType(params.Sort.ByNumber),
		)
		sq.Add(q, *params.Pagination.AfterID)
	}

	sq.WhereOptional(func() { sq.AddCondition(params) })
	sq.SortByNumber(params.Sort.ByNumber)
	sq.Add("LIMIT ?", params.Pagination.Limit)

	query, args := sq.Make()

	err = pgxscan.Select(ctx, r.db.TxOrDB(ctx), &tasks, query, args...)
	if err != nil {
		return tList, errors.Wrap(err, "GetAllByParams Task repo:")
	}

	tList.Tasks = tasks
	if len(tasks) != 0 {
		tList.Pagination.AfterID = tasks[len(tasks)-1].ID
	}

	return tList, nil
}
