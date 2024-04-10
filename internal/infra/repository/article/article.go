package article

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

func (r Repository) Create(ctx context.Context, dto domain.ArticleCreateInput) (a domain.Article, err error) {
	sq := sql_query_maker.NewQueryMaker(4)

	sq.Add(
		`
	INSERT INTO article (author_id, title, content, categories)
	VALUES (?, ?, ?, ?)
	RETURNING id, author_id, title, content, categories, created_at
	`,
		dto.AuthorID, dto.Title, dto.Content, dto.Categories,
	)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &a, query, args...)
	if err == nil {
		return a, nil
	}

	var pgError *pgconn.PgError
	if ok := errors.As(err, &pgError); !ok {
		return a, errors.Wrap(err, "Create Article repo:")
	}

	switch pgError.Code {
	case postgres.ERRCODE_UNIQUE_VIOLATION:
		err = &struct_errors.ErrExist{Err: err, Msg: "Title already exist"}
	}

	return a, errors.Wrap(err, "Create Article repo:")
}

func (r Repository) Update(ctx context.Context, dto domain.ArticleUpdateInput) (a domain.Article, err error) {
	sq := sql_query_maker.NewQueryMaker(4)

	sq.Add("UPDATE article SET")

	if dto.Title != nil {
		sq.Add("title = ?,", *dto.Title)
	}

	if dto.Content != nil {
		sq.Add("content = ?,", *dto.Content)
	}

	if dto.Categories != nil {
		sq.Add("categories =?,", dto.Categories)
	}

	sq.Where("id = ?", dto.ID)
	sq.Add("RETURNING id, author_id, title, content, categories, created_at")

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &a, query, args...)
	if err != nil {
		return a, errors.Wrap(err, "Update Article repo:")
	}

	return a, nil
}

func (r Repository) Delete(ctx context.Context, id string) error {
	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add("DELETE FROM article WHERE id = ?", id)

	query, args := sq.Make()

	res, err := r.db.TxOrDB(ctx).Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "Delete Article repo:")
	}

	if res.RowsAffected() == 0 {
		err = errors.New("Article not found!")

		return errors.Wrap(err, "Delete Article repo:")
	}

	return nil
}

func (r Repository) GetByID(ctx context.Context, id string) (a domain.Article, err error) {
	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add(
		`
	SELECT id, author_id, title, content, categories, created_at
	FROM article
	WHERE id = ?
	`,
		id)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &a, query, args...)
	if err != nil {
		return a, errors.Wrap(err, "GetByID Article repo:")
	}

	return a, nil
}

func (r Repository) GetPracticeArticle(ctx context.Context) (a domain.Article, err error) {
	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add(
		`
	SELECT id, author_id, title, content, categories, created_at
	FROM article
	WHERE title = ?
	`,
		domain.PracticeArticleName)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &a, query, args...)
	if err != nil {
		return a, errors.Wrap(err, "GetPracticeArticle Article repo:")
	}

	return a, nil
}

func (r Repository) GetAllByParams(ctx context.Context, params domain.ArticleParams) (aList domain.ArticleList, err error) {
	articles := []domain.Article{}
	sq := newFilter(r.cfg, 15)

	sq.Add("SELECT id, author_id, title, categories, created_at FROM article a")

	if params.Pagination.AfterID != nil {
		q := fmt.Sprintf(
			"WHERE a.title != ? AND created_at %s (SELECT created_at FROM article WHERE id = ?)",
			db.GetLetterGreaterOrLessBySortType(params.Sort.ByDate),
		)
		sq.Add(q, domain.PracticeArticleName, *params.Pagination.AfterID)
	} else {
		sq.Add("WHERE a.title != ?", domain.PracticeArticleName)
	}

	sq.AddCondition(params)
	sq.SortByCreatedAt(params.Sort.ByDate)
	sq.Add("LIMIT ?", params.Pagination.Limit)

	query, args := sq.Make()

	err = pgxscan.Select(ctx, r.db.TxOrDB(ctx), &articles, query, args...)
	if err != nil {
		return domain.ArticleList{}, errors.Wrap(err, "GetAllByParams Article repo:")
	}

	aList.Articles = articles
	if len(articles) != 0 {
		aList.Pagination.AfterID = articles[len(articles)-1].ID
	}

	return aList, nil
}

func (r Repository) GetAvailableAttributes(ctx context.Context) (domain.ArticleAttributes, error) {
	categories := []string{}
	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add(
		`
	SELECT DISTINCT unnest(categories)
	FROM article a
	WHERE a.title != ?
	`,
		domain.PracticeArticleName)

	query, args := sq.Make()

	err := pgxscan.Select(ctx, r.db.TxOrDB(ctx), &categories, query, args...)
	if err != nil {
		return domain.ArticleAttributes{}, errors.Wrap(err, "GetAvailableAttributes Article repo:")
	}

	return domain.ArticleAttributes{Categories: categories}, nil
}
