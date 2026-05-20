package transaction

import (
	"context"
	"database/sql"
)

// DBTX abstracts *sql.DB and *sql.Tx so repositories work inside or outside transactions.
type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type txCtxKey struct{}

type txHolder struct {
	tx *sql.Tx
}

func txFromContext(ctx context.Context) *txHolder {
	h, _ := ctx.Value(txCtxKey{}).(*txHolder)
	return h
}

// NewDBGetter returns a function that resolves the active DBTX from context,
// falling back to db when no transaction is in progress.
func NewDBGetter(db *sql.DB) func(ctx context.Context) DBTX {
	return func(ctx context.Context) DBTX {
		if h := txFromContext(ctx); h != nil {
			return h.tx
		}
		return db
	}
}

type pgTransactor struct {
	db *sql.DB
}

func NewTransactor(db *sql.DB) Transactor {
	return &pgTransactor{db: db}
}

func (t *pgTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if txFromContext(ctx) != nil {
		return fn(ctx)
	}

	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	txCtx := context.WithValue(ctx, txCtxKey{}, &txHolder{tx: tx})
	if err := fn(txCtx); err != nil {
		return err
	}
	return tx.Commit()
}
