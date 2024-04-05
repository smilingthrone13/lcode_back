package task_template

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

func (r *Repository) Create(ctx context.Context, dto domain.TaskTemplateCreate) (tt domain.TaskTemplate, err error) {
	sq := sql_query_maker.NewQueryMaker(4)

	sq.Add(
		`
	INSERT INTO task_template(task_id, language_id, template, wrapper)
	VALUES (?, ?, ?, ?)
	RETURNING id, task_id, language_id, template, wrapper
	`,
		dto.TaskID, dto.LanguageID, dto.Template, dto.Wrapper,
	)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &tt, query, args...)
	if err != nil {
		return tt, errors.Wrap(err, "Create TaskTemplate repo:")
	}

	return tt, nil
}

func (r *Repository) Update(
	ctx context.Context,
	id string,
	dto domain.TaskTemplateUpdate,
) (tt domain.TaskTemplate, err error) {
	sq := sql_query_maker.NewQueryMaker(3)

	sq.Add("UPDATE task_template SET")

	if dto.Template != nil {
		sq.Add("template = ?,", *dto.Template)
	}

	if dto.Wrapper != nil {
		sq.Add("wrapper = ?,", *dto.Wrapper)
	}

	sq.Where("id = (SELECT id FROM task_template WHERE id = ? FOR UPDATE)", id)
	sq.Add("RETURNING id, task_id, language_id, template, wrapper")

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &tt, query, args...)
	if err != nil {
		return tt, errors.Wrap(err, "Update TaskTemplate repo:")
	}

	return tt, nil
}

func (r *Repository) Delete(ctx context.Context, id string) (tt domain.TaskTemplate, err error) {
	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add(
		`
	DELETE FROM task_template WHERE id = ?
	RETURNING id, task_id, language_id, template, wrapper
	`,
		id)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &tt, query, args...)
	if err != nil {
		return tt, errors.Wrap(err, "Delete TaskTemplate repo:")
	}

	return tt, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (tt domain.TaskTemplate, err error) {
	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add(
		`
	SELECT id, task_id, language_id, template, wrapper
	FROM task_template WHERE id = ?
	`,
		id)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &tt, query, args...)
	if err != nil {
		return tt, errors.Wrap(err, "GetByID TaskTemplate repo:")
	}

	return tt, nil
}

func (r *Repository) GetAllByTaskID(ctx context.Context, id string) (tt []domain.TaskTemplate, err error) {
	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add(
		`
	SELECT id, task_id, language_id, template, wrapper
	FROM task_template WHERE task_id = ?
	`,
		id)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &tt, query, args...)
	if err != nil {
		return tt, errors.Wrap(err, "GetAllByTaskID TaskTemplate repo:")
	}

	return tt, nil
}
