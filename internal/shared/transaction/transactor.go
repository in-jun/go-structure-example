package transaction

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type IsolationLevel int

const (
	Optimistic  IsolationLevel = iota // Serializable + auto-retry on serialization failure
	Pessimistic                       // ReadCommitted; callers must use SELECT FOR UPDATE
)

type TxOption func(*txConfig)

type txConfig struct {
	isolation IsolationLevel
}

func WithIsolation(level IsolationLevel) TxOption {
	return func(c *txConfig) {
		c.isolation = level
	}
}

type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error, opts ...TxOption) error
}

type txCtxKey struct{}
type postCommitKey struct{}

func txFromContext(ctx context.Context) *sql.Tx {
	tx, ok := ctx.Value(txCtxKey{}).(*sql.Tx)
	if !ok {
		return nil
	}
	return tx
}

func RegisterPostCommit(ctx context.Context, fn func()) {
	if hooks, ok := ctx.Value(postCommitKey{}).(*[]func()); ok {
		*hooks = append(*hooks, fn)
	} else {
		fn()
	}
}

func NewDBGetter(db *sql.DB) func(ctx context.Context) DBTX {
	return func(ctx context.Context) DBTX {
		if tx := txFromContext(ctx); tx != nil {
			return tx
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

const maxSerializationRetries = 3

func isRetryable(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "40001" || pgErr.Code == "40P01"
	}
	return false
}

func (t *pgTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error, opts ...TxOption) error {
	if txFromContext(ctx) != nil {
		return fn(ctx)
	}

	cfg := txConfig{isolation: Optimistic}
	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.isolation == Pessimistic {
		return t.executeTx(ctx, sql.LevelReadCommitted, fn)
	}

	var lastErr error
	for attempt := 0; attempt <= maxSerializationRetries; attempt++ {
		if attempt > 0 && ctx.Err() != nil {
			return ctx.Err()
		}

		lastErr = t.executeTx(ctx, sql.LevelSerializable, fn)
		if lastErr == nil || !isRetryable(lastErr) {
			return lastErr
		}
	}
	return lastErr
}

func (t *pgTransactor) executeTx(ctx context.Context, level sql.IsolationLevel, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, &sql.TxOptions{Isolation: level})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var hooks []func()
	txCtx := context.WithValue(ctx, txCtxKey{}, tx)
	txCtx = context.WithValue(txCtx, postCommitKey{}, &hooks)

	if err := fn(txCtx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	for _, hook := range hooks {
		hook()
	}
	return nil
}
