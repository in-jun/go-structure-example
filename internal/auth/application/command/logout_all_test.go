package command

import (
	"context"
	"testing"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

func TestLogoutAllHandler_Success(t *testing.T) {
	h := NewLogoutAllHandler(&mockTokenRepo{}, &mockTokenGen{})

	err := h.Handle(context.Background(), LogoutAll{UserID: testUUID})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestLogoutAllHandler_EmptyUserID(t *testing.T) {
	h := NewLogoutAllHandler(&mockTokenRepo{}, &mockTokenGen{})

	err := h.Handle(context.Background(), LogoutAll{UserID: ""})
	if err == nil {
		t.Fatal("expected error for empty userID, got nil")
	}
}

func TestLogoutAllHandler_RepositoryError(t *testing.T) {
	h := NewLogoutAllHandler(&mockTokenRepo{err: errors.Internal("db error")}, &mockTokenGen{})

	err := h.Handle(context.Background(), LogoutAll{UserID: testUUID})
	if err == nil {
		t.Fatal("expected error for repository failure, got nil")
	}
}
