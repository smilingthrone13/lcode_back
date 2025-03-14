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

func (r *Repository) CreateDefault(ctx context.Context, user domain.User) error {
	aText := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus.\nSuspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor.\nCras elementum ultrices diam. Maecenas ligula massa, varius a, semper sagittis, dapibus gravida, tellus.\nNulla vitae elit. Nulla facilisi. Ut fringilla. Suspendisse eu ligula. Etiam porta sem."
	aTitle := "Practice Article"
	aCategories := []string{"Practice"}
	sq := sql_query_maker.NewQueryMaker(5)

	sq.Add(
		`
	INSERT INTO article (id, author_id, title, content, categories)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT ON CONSTRAINT articles_pk DO UPDATE SET author_id = excluded.author_id, 
	                          title = excluded.title, 
	                          content = excluded.content, 
	                          categories = excluded.categories
	`,
		domain.PracticeArticleID, user.ID, aTitle, aText, aCategories,
	)

	query, args := sq.Make()

	_, err := r.db.TxOrDB(ctx).Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "CreateDefault Article repo:")
	}

	return nil
}

func (r *Repository) Create(ctx context.Context, inp domain.ArticleCreateInput) (a domain.Article, err error) {
	var id string
	sq := sql_query_maker.NewQueryMaker(4)

	sq.Add(
		`
	INSERT INTO article (author_id, title, content, categories)
	VALUES (?, ?, ?, ?)
	RETURNING id
	`,
		inp.AuthorID, inp.Title, inp.Content, inp.Categories,
	)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &id, query, args...)
	if err != nil {
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

	a, err = r.GetByID(ctx, id)
	if err != nil {
		return a, errors.Wrap(err, "Create Article repo:")
	}

	return a, nil
}

func (r *Repository) Update(ctx context.Context, inp domain.ArticleUpdateInput) (a domain.Article, err error) {
	sq := sql_query_maker.NewQueryMaker(4)

	sq.Add("UPDATE article SET")

	if inp.Title != nil {
		sq.Add("title = ?,", *inp.Title)
	}

	if inp.Content != nil {
		sq.Add("content = ?,", *inp.Content)
	}

	if inp.Categories != nil {
		sq.Add("categories = ?,", inp.Categories)
	}

	sq.Where("id = ?", inp.ID)

	query, args := sq.Make()

	res, err := r.db.TxOrDB(ctx).Exec(ctx, query, args...)
	if err != nil {
		var pgError *pgconn.PgError
		if ok := errors.As(err, &pgError); !ok {
			return a, errors.Wrap(err, "Update Article repo:")
		}

		switch pgError.Code {
		case postgres.ERRCODE_UNIQUE_VIOLATION:
			err = &struct_errors.ErrExist{Err: err, Msg: "Title already exist"}
		}

		return a, errors.Wrap(err, "Create Article repo:")
	}

	if res.RowsAffected() == 0 {
		err = errors.New("Article not found!")

		return a, errors.Wrap(err, "Update Article repo:")
	}

	a, err = r.GetByID(ctx, inp.ID)
	if err != nil {
		return a, errors.Wrap(err, "Update Article repo:")
	}

	return a, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
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

func (r *Repository) GetByID(ctx context.Context, id string) (a domain.Article, err error) {
	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add(
		`
	SELECT 
	    a.id AS id, title, content, categories, created_at,
	    u.id AS user_id, u.username AS username, u.first_name AS first_name, u.last_name AS last_name
	FROM article a 
	    JOIN "user" u ON a.author_id = u.id
	WHERE a.id = ?
	`,
		id)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &a, query, args...)
	if err != nil {
		return a, errors.Wrap(err, "GetByID Article repo:")
	}

	return a, nil
}

func (r *Repository) GetAllByParams(ctx context.Context, params domain.ArticleParams) (aList domain.ArticleList, err error) {
	articles := []domain.Article{}
	sq := newFilter(r.cfg, 15)

	sq.Add(
		`
	SELECT 
	    a.id AS id, title, content, categories, created_at,
	    u.id AS user_id, u.username AS username, u.first_name AS first_name, u.last_name AS last_name
	FROM article a 
	    JOIN "user" u ON a.author_id = u.id
	`,
	)

	if params.Pagination.AfterID != nil {
		q := fmt.Sprintf(
			"WHERE a.id != ? AND a.created_at %s (SELECT created_at FROM article WHERE id = ?)",
			db.GetLetterGreaterOrLessBySortType(params.Sort.ByDate),
		)
		sq.Add(q, domain.PracticeArticleID, *params.Pagination.AfterID)
	} else {
		sq.Add("WHERE a.id != ?", domain.PracticeArticleID)
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

func (r *Repository) GetAvailableAttributes(ctx context.Context) (domain.ArticleAttributes, error) {
	categories := []string{}
	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add(
		`
	SELECT DISTINCT unnest(categories)
	FROM article a
	WHERE a.id != ?
	`,
		domain.PracticeArticleID)

	query, args := sq.Make()

	err := pgxscan.Select(ctx, r.db.TxOrDB(ctx), &categories, query, args...)
	if err != nil {
		return domain.ArticleAttributes{}, errors.Wrap(err, "GetAvailableAttributes Article repo:")
	}

	return domain.ArticleAttributes{Categories: categories}, nil
}
