package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

const (
	ERRCODE_FOREIGN_KEY_VIOLATION = "23503"
	ERRCODE_UNIQUE_VIOLATION      = "23505"
)

func NewDBManger(db *pgxpool.Pool) *DbManager {
	return &DbManager{
		db: db,
	}
}

type DbManager struct {
	db *pgxpool.Pool
}

type TxKey struct{}

func (d *DbManager) TxOrDB(ctx context.Context) QueryExecutor {
	if tx, ok := ctx.Value(TxKey{}).(*Tx); ok {
		return tx.tx
	}

	return d.db
}

func (d *DbManager) GetDb() *pgxpool.Pool {
	return d.db
}

func NewTransactionProvider(db *pgxpool.Pool) *TransactionProvider {
	t := &TransactionProvider{db: db, txQueue: make(chan *Tx, 2000)}

	go t.afterSuccessWorker()

	return t
}

type TransactionProvider struct {
	db      *pgxpool.Pool
	txQueue chan *Tx
}

func (tm *TransactionProvider) afterSuccessWorker() {
	for tx := range tm.txQueue {
		for i := 0; i < len(tx.afterSuccessQueue); i++ {
			tx.afterSuccessQueue[i]()
		}
	}
}

func (tm *TransactionProvider) GetTxForParticipant(ctx context.Context) (TxForParticipant, error) {
	if tx, ok := ctx.Value(TxKey{}).(*Tx); ok {
		return tx, nil
	}

	return nil, errors.New("transaction not found in ctx")
}

func (tm *TransactionProvider) NewTx(ctx context.Context, opts *pgx.TxOptions) (Transaction, error) {
	var tx pgx.Tx
	var err error

	if opts == nil {
		tx, err = tm.db.Begin(ctx)
	} else {
		tx, err = tm.db.BeginTx(ctx, *opts)
	}

	if err != nil {
		err = errors.Wrap(err, "NewTx postgres pkg")

		return nil, err
	}

	return &Tx{tx: tx, txQueue: tm.txQueue}, nil
}

type Tx struct {
	tx                pgx.Tx
	afterSuccessQueue []func()
	txQueue           chan<- *Tx
}

func (t *Tx) Commit(ctx context.Context) error {
	if err := t.tx.Commit(ctx); err != nil {
		return errors.Wrap(err, "Commit postgres pkg")
	}

	t.txQueue <- t

	return nil
}

func (t *Tx) AfterSuccess(ctx context.Context, f func()) {
	t.afterSuccessQueue = append(t.afterSuccessQueue, f)
}

func (t *Tx) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

type Transaction interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	AfterSuccess(ctx context.Context, f func())
}

type TxForParticipant interface {
	AfterSuccess(ctx context.Context, f func())
}

type QueryExecutor interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
}

type errConstructor func(err error) error

type ErrPair struct {
	code           string
	errConstructor errConstructor
}

func Error(code string, err errConstructor) ErrPair {
	return ErrPair{code: code, errConstructor: err}
}
