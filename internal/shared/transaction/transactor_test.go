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
	tx := txFromContext(context.Background())
	if tx != nil {
		t.Error("expected nil tx from empty context")
	}
}

func TestWithinTransaction_NestedReusesExistingTx(t *testing.T) {
	fakeTx := new(sql.Tx)
	txCtx := context.WithValue(context.Background(), txCtxKey{}, fakeTx)

	called := 0
	transactor := &pgTransactor{db: nil}

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

	fakeTx := new(sql.Tx)
	txCtx := context.WithValue(context.Background(), txCtxKey{}, fakeTx)
	result := getter(txCtx)
	if result != fakeTx {
		t.Errorf("expected tx from context, got %v", result)
	}
}

func TestRegisterPostCommit_NoTransaction(t *testing.T) {
	called := false
	RegisterPostCommit(context.Background(), func() { called = true })
	if !called {
		t.Error("expected hook to be called immediately when no transaction is in progress")
	}
}
