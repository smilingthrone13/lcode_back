package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"lcode/config"
	"lcode/pkg/postgres"
	"log"
	"time"
)

func NewPostgresDB(cfg *config.Config) (db *postgres.DbManager, err error) {
	ctx := context.Background()

	c, err := pgxpool.ParseConfig(cfg.DBConfig.Path)
	if err != nil {
		log.Fatal(err)
	}

	c.MaxConns = 50

	c.MinConns = 2

	c.HealthCheckPeriod = time.Second * 1

	dbPg, err := pgxpool.NewWithConfig(ctx, c)
	if err != nil {
		log.Fatal(errors.Wrap(err, "NewPostgresDB database:"))
	}

	db = postgres.NewDBManger(dbPg)

	return db, nil
}
