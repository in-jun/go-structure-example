package command

import (
	"context"
	"testing"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

func TestDeleteHandler_Success(t *testing.T) {
	h := NewDeleteHandler(&mockUserRepo{user: makeUser()})

	err := h.Handle(context.Background(), Delete{UserID: 1})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestDeleteHandler_RepositoryError(t *testing.T) {
	h := NewDeleteHandler(&mockUserRepo{err: errors.Internal("db error")})

	err := h.Handle(context.Background(), Delete{UserID: 1})
	if err == nil {
		t.Fatal("expected error for repository failure, got nil")
	}
}
