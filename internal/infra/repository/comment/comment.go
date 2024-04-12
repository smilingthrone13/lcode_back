package comment

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
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

func (r *Repository) Create(ctx context.Context, dto domain.CommentCreateDTO) (c domain.Comment, err error) {
	var id string
	sq := sql_query_maker.NewQueryMaker(4)

	q := fmt.Sprintf(
		`
	INSERT INTO %s (parent_id, entity_id, author_id, comment_text)
	VALUES (?, ?, ?, ?)
	RETURNING id
	`,
		dto.OriginType)

	sq.Add(q, dto.Input.ParentID, dto.Input.EntityID, dto.Input.AuthorID, dto.Input.Text)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &id, query, args...)
	if err != nil {
		return c, errors.Wrap(err, "Create Comment repo:")
	}

	c, err = r.getByID(ctx, dto.OriginType, id)
	if err != nil {
		return c, errors.Wrap(err, "Create Comment repo:")
	}

	return c, nil
}

func (r *Repository) Update(ctx context.Context, dto domain.CommentUpdateDTO) (c domain.Comment, err error) {
	sq := sql_query_maker.NewQueryMaker(4)

	c, err = r.getByID(ctx, dto.OriginType, dto.Input.ID)
	if err != nil {
		return c, errors.Wrap(err, "Update Comment repo:")
	}

	if c.Author.UserID != dto.User.ID && !dto.User.IsAdmin {
		err = struct_errors.NewForbiddenErr(fmt.Errorf("no access rights"))

		return c, errors.Wrap(err, "Update Comment repo:")
	}

	sq.Add(fmt.Sprintf("UPDATE %s SET", dto.OriginType))

	if dto.Input.Text != nil {
		sq.Add("comment_text = ?", *dto.Input.Text)
	}

	sq.Where("id = ?", dto.Input.ID)

	query, args := sq.Make()

	res, err := r.db.TxOrDB(ctx).Exec(ctx, query, args...)
	if err != nil {
		return c, errors.Wrap(err, "Update Comment repo:")
	}

	if res.RowsAffected() == 0 {
		err = errors.New("Comment not found!")

		return c, errors.Wrap(err, "Update Comment repo:")
	}

	c, err = r.getByID(ctx, dto.OriginType, dto.Input.ID)
	if err != nil {
		return c, errors.Wrap(err, "Update Comment repo:")
	}

	return c, nil
}

func (r *Repository) Delete(ctx context.Context, dto domain.CommentDeleteDTO) error {
	sq := sql_query_maker.NewQueryMaker(2)

	c, err := r.getByID(ctx, dto.OriginType, dto.ID)
	if err != nil {
		return errors.Wrap(err, "Delete Comment repo:")
	}

	if c.Author.UserID != dto.User.ID && !dto.User.IsAdmin {
		err = struct_errors.NewForbiddenErr(fmt.Errorf("no access rights"))

		return errors.Wrap(err, "Delete Comment repo:")
	}

	sq.Add(fmt.Sprintf("DELETE FROM %s WHERE id = ?", dto.OriginType), dto.ID)

	query, args := sq.Make()

	res, err := r.db.TxOrDB(ctx).Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "Delete Comment repo:")
	}

	if res.RowsAffected() == 0 {
		err = errors.New("Comment not found!")

		return errors.Wrap(err, "Delete Comment repo:")
	}

	return nil
}

func (r *Repository) GetThreadsByParamsAndEntityID(
	ctx context.Context,
	dto domain.CommentParamsDTO,
) (tl domain.ThreadList, err error) {
	threadHeads, err := r.getThreadHeads(ctx, dto.OriginType, dto.EntityID, dto.Input)
	if err != nil {
		return tl, errors.Wrap(err, "GetThreadsByParamsAndEntityID Comment repo:")
	}

	headIDs := make([]string, 0, len(threadHeads))
	for i := range threadHeads {
		headIDs = append(headIDs, threadHeads[i].ID)
	}

	replies, err := r.getRepliesByCommentIDs(ctx, dto.OriginType, headIDs)
	if err != nil {
		return tl, errors.Wrap(err, "GetThreadsByParamsAndEntityID Comment repo:")
	}

	tl.Threads = r.splitIntoThreads(threadHeads, replies)
	if len(tl.Threads) != 0 {
		tl.Pagination.AfterID = tl.Threads[len(tl.Threads)-1].Comment.ID
	}

	return tl, nil
}

func (r *Repository) getByID(ctx context.Context, origin domain.CommentOriginType, id string) (c domain.Comment, err error) {
	var comms []domain.Comment
	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add(
		fmt.Sprintf(
			`
			SELECT 
			    c.id AS id, parent_id, entity_id, comment_text, created_at,
			    u.id AS user_id, u.username AS username, u.first_name AS first_name, u.last_name AS last_name
			FROM %s c 
			    JOIN "user" u ON c.author_id = u.id
			WHERE c.id = ?
			`,
			origin,
		),
		id)

	query, args := sq.Make()

	err = pgxscan.Select(ctx, r.db.TxOrDB(ctx), &comms, query, args...)
	if err != nil {
		return c, errors.Wrap(err, "getByID Comment repo:")
	}

	if len(comms) < 1 {
		err = struct_errors.NewErrNotFound("Comment not found", nil)

		return c, errors.Wrap(err, "getByID Comment repo:")
	}

	return comms[0], nil
}

func (r *Repository) getThreadHeads(
	ctx context.Context,
	origin domain.CommentOriginType,
	entityID string,
	params domain.CommentParamsInput,
) (heads []domain.Comment, err error) {
	sq := newFilter(r.cfg, 15)

	sq.Add(
		fmt.Sprintf(
			`
			SELECT 
			    c.id AS id, parent_id, entity_id, comment_text, created_at,
			    u.id AS user_id, u.username AS username, u.first_name AS first_name, u.last_name AS last_name
			FROM %s c 
			    JOIN "user" u ON c.author_id = u.id
			`,
			origin,
		),
	)

	if params.Pagination.AfterID != nil {
		q := fmt.Sprintf(
			"WHERE c.entity_id = ? AND c.parent_id IS NULL AND c.created_at %s (SELECT created_at FROM %s WHERE id = ?)",
			db.GetLetterGreaterOrLessBySortType(params.Sort.ByDate), origin,
		)
		sq.Add(q, entityID, *params.Pagination.AfterID)
	} else {
		sq.Add("WHERE c.entity_id = ? AND c.parent_id IS NULL", entityID)
	}

	sq.SortByCreatedAt(params.Sort.ByDate)
	sq.Add("LIMIT ?", params.Pagination.Limit)

	query, args := sq.Make()

	err = pgxscan.Select(ctx, r.db.TxOrDB(ctx), &heads, query, args...)
	if err != nil {
		return heads, errors.Wrap(err, "getThreadHeads Article repo:")
	}

	return heads, nil
}

func (r *Repository) getRepliesByCommentIDs(
	ctx context.Context,
	origin domain.CommentOriginType,
	commentIDs []string,
) (replies []domain.Comment, err error) {
	sq := sql_query_maker.NewQueryMaker(1)

	sq.Add(
		fmt.Sprintf(
			`
			SELECT 
			    c.id AS id, parent_id, entity_id, comment_text, created_at,
			    u.id AS user_id, u.username AS username, u.first_name AS first_name, u.last_name AS last_name
			FROM %s c 
			    JOIN "user" u ON c.author_id = u.id
			WHERE c.parent_id = ANY(?)
			ORDER BY c.created_at DESC
			`,
			origin,
		),
		commentIDs,
	)

	query, args := sq.Make()

	err = pgxscan.Select(ctx, r.db.TxOrDB(ctx), &replies, query, args...)
	if err != nil {
		return replies, errors.Wrap(err, "getRepliesByCommentIDs Article repo:")
	}

	return replies, nil
}

func (r *Repository) splitIntoThreads(threadHeads []domain.Comment, replies []domain.Comment) []domain.Thread {
	repliesMap := make(map[string][]domain.Comment, len(threadHeads))

	for i := range replies {
		_, ok := repliesMap[*replies[i].ParentID]
		if !ok {
			repliesMap[*replies[i].ParentID] = []domain.Comment{replies[i]}
			continue
		}

		repliesMap[*replies[i].ParentID] = append(
			repliesMap[*replies[i].ParentID],
			replies[i],
		)
	}

	threads := make([]domain.Thread, 0, len(threadHeads))
	for i := range threadHeads {
		thisThreadHead := threadHeads[i]
		thisThreadReplies, ok := repliesMap[thisThreadHead.ID]
		if !ok {
			thisThreadReplies = []domain.Comment{}
		}

		threads = append(
			threads,
			domain.Thread{
				Comment: thisThreadHead,
				Replies: thisThreadReplies,
			},
		)
	}

	return threads
}
