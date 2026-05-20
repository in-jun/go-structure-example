package transaction

import (
	"context"
	"database/sql"
	"testing"
)

func TestNewDBGetter_ReturnsFunction(t *testing.T) {
	db := &sql.DB{}
	getter := NewDBGetter(db)
	if getter == nil {
		t.Error("expected non-nil getter function")
	}
}

func TestTxFromContext_Empty(t *testing.T) {
	h := txFromContext(context.Background())
	if h != nil {
		t.Error("expected nil holder from empty context")
	}
}

func TestWithinTransaction_NestedReusesExistingTx(t *testing.T) {
	// Inject a txHolder into context to simulate an active transaction.
	txCtx := context.WithValue(context.Background(), txCtxKey{}, &txHolder{tx: nil})

	called := 0
	transactor := &mysqlTransactor{db: nil}

	err := transactor.WithinTransaction(txCtx, func(ctx context.Context) error {
		called++
		if ctx != txCtx {
			t.Error("nested transaction should receive the same context")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called != 1 {
		t.Errorf("expected fn called once, got %d", called)
	}
}

func TestNewDBGetter_UsesTxFromContext(t *testing.T) {
	db := &sql.DB{}
	getter := NewDBGetter(db)

	txCtx := context.WithValue(context.Background(), txCtxKey{}, &txHolder{tx: nil})
	result := getter(txCtx)
	// When a txHolder is in context, DBTX should be the tx (nil here), not the db.
	if result != (*sql.Tx)(nil) {
		t.Errorf("expected tx from context, got %v", result)
	}
}
