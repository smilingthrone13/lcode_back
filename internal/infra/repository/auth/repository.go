package auth

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

func (r *Repository) CreateUser(ctx context.Context, dto domain.CreateUserEntity) (user domain.User, err error) {
	sq := sql_query_maker.NewQueryMaker(4)

	sq.Add(
		`INSERT INTO "user" (first_name, last_name, username, password_hash) 
			   VALUES (?, ?, ?, ?) 
               RETURNING id, first_name, last_name, username, password_hash, is_admin`,
		dto.FirstName,
		dto.LastName,
		dto.Username,
		dto.PasswordHash,
	)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &user, query, args...)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "CreateUser auth repo")
	}

	return user, nil
}

func (r *Repository) ChangeUserAdminStatus(ctx context.Context, dto domain.ChangeUserAdminPermissionDTO) error {
	sq := sql_query_maker.NewQueryMaker(2)

	sq.Add(
		`UPDATE "user" SET is_admin = ? WHERE id = ?`,
		dto.IsAdmin,
		dto.UserID,
	)

	query, args := sq.Make()

	_, err := r.db.TxOrDB(ctx).Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "ChangeUserAdminStatus auth repo")
	}

	return nil
}

func (r *Repository) UserByID(ctx context.Context, id string) (user domain.User, err error) {
	sq := sql_query_maker.NewQueryMaker(2)

	sq.Add(
		`SELECT id, first_name, last_name, username, password_hash, is_admin FROM "user" WHERE id =?`,
		id,
	)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &user, query, args...)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "UserByID auth repo")
	}

	return user, nil
}

func (r *Repository) Users(ctx context.Context) ([]domain.User, error) {
	sq := sql_query_maker.NewQueryMaker(0)

	users := []domain.User{}

	sq.Add(`SELECT id, first_name, last_name, username, password_hash, is_admin FROM "user"`)

	query, args := sq.Make()

	err := pgxscan.Select(ctx, r.db.TxOrDB(ctx), &users, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Users auth repo")
	}

	return users, nil
}

func (r *Repository) UserByUsername(ctx context.Context, username string) (user domain.User, err error) {
	sq := sql_query_maker.NewQueryMaker(2)

	sq.Add(
		`SELECT id, first_name, last_name, username, password_hash, is_admin FROM "user" WHERE username = ?`,
		username,
	)

	query, args := sq.Make()

	err = pgxscan.Get(ctx, r.db.TxOrDB(ctx), &user, query, args...)
	if err != nil {
		return domain.User{}, errors.Wrap(err, "UserByUsername auth repo")
	}

	return user, nil
}
