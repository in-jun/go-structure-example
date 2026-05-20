package command

import (
	"context"
	"testing"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

func TestLogoutHandler_Success(t *testing.T) {
	h := NewLogoutHandler(&mockTokenRepo{}, &mockTokenGen{})

	err := h.Handle(context.Background(), Logout{RefreshToken: "some-refresh-token", AccessTokenJTI: "jti-123"})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestLogoutHandler_EmptyToken(t *testing.T) {
	h := NewLogoutHandler(&mockTokenRepo{}, &mockTokenGen{})

	err := h.Handle(context.Background(), Logout{RefreshToken: ""})
	if err == nil {
		t.Fatal("expected error for empty refresh token, got nil")
	}
}

func TestLogoutHandler_RepositoryError(t *testing.T) {
	h := NewLogoutHandler(&mockTokenRepo{err: errors.Internal("db error")}, &mockTokenGen{})

	err := h.Handle(context.Background(), Logout{RefreshToken: "some-token"})
	if err == nil {
		t.Fatal("expected error for repository failure, got nil")
	}
}

func TestLogoutHandler_WithoutJTI(t *testing.T) {
	h := NewLogoutHandler(&mockTokenRepo{}, &mockTokenGen{})

	err := h.Handle(context.Background(), Logout{RefreshToken: "some-refresh-token"})
	if err != nil {
		t.Fatalf("Handle() without JTI should succeed, got %v", err)
	}
}
