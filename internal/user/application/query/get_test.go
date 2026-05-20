package query

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/user/domain"
	"github.com/in-jun/go-structure-example/internal/user/domain/entity"
)

type mockUserRepo struct {
	user *entity.User
	err  error
}

func (m *mockUserRepo) FindByID(_ context.Context, _ uint) (*entity.User, error) {
	return m.user, m.err
}
func (m *mockUserRepo) Update(_ context.Context, _ *entity.User) error { return m.err }
func (m *mockUserRepo) Delete(_ context.Context, _ uint) error         { return m.err }

var _ domain.UserRepository = (*mockUserRepo)(nil)

func TestGetHandler_Success(t *testing.T) {
	now := time.Now()
	u, _ := entity.ReconstructUser(1, "test@example.com", "hashed", "Test User", now, now)
	h := NewGetHandler(&mockUserRepo{user: u})

	result, err := h.Handle(context.Background(), Get{UserID: 1})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if result.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", result.Email, "test@example.com")
	}
	if result.Name != "Test User" {
		t.Errorf("Name = %q, want %q", result.Name, "Test User")
	}
}

func TestGetHandler_NotFound(t *testing.T) {
	h := NewGetHandler(&mockUserRepo{err: errors.NotFound("user not found")})

	_, err := h.Handle(context.Background(), Get{UserID: 99})
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}
