package command

import (
	"context"
	"testing"

	"github.com/in-jun/go-structure-example/internal/auth/domain"
	"github.com/in-jun/go-structure-example/internal/auth/domain/entity"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

type mockUserRepo struct {
	user *entity.User
	err  error
}

func (m *mockUserRepo) Save(_ context.Context, u *entity.User) error {
	u.SetID(1)
	return m.err
}
func (m *mockUserRepo) FindByEmail(_ context.Context, _ string) (*entity.User, error) {
	return m.user, m.err
}

type mockHasher struct{}

func (m *mockHasher) Hash(password string) (string, error) { return "hashed_" + password, nil }
func (m *mockHasher) Compare(hashed, plain string) bool    { return hashed == "hashed_"+plain }

type noopTransactor struct{}

func (n *noopTransactor) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

var _ domain.UserRepository = (*mockUserRepo)(nil)
var _ domain.PasswordHasher = (*mockHasher)(nil)
var _ transaction.Transactor = (*noopTransactor)(nil)

func TestRegisterHandler_Success(t *testing.T) {
	h := NewRegisterHandler(&mockUserRepo{}, &mockHasher{}, &noopTransactor{})
	err := h.Handle(context.Background(), Register{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestRegisterHandler_InvalidEmail(t *testing.T) {
	h := NewRegisterHandler(&mockUserRepo{}, &mockHasher{}, &noopTransactor{})
	err := h.Handle(context.Background(), Register{
		Email:    "not-an-email",
		Password: "password123",
		Name:     "Test",
	})
	if err == nil {
		t.Fatal("expected error for invalid email")
	}
}

func TestRegisterHandler_DuplicateEmail(t *testing.T) {
	existing, _ := entity.NewUser("test@example.com", "hashed", "Existing")
	existing.SetID(1)
	h := NewRegisterHandler(&mockUserRepo{user: existing}, &mockHasher{}, &noopTransactor{})
	err := h.Handle(context.Background(), Register{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test",
	})
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
	var ce errors.CustomError
	if !asCustomError(err, &ce) || ce.Status != 409 {
		t.Errorf("expected 409 Conflict, got %v", err)
	}
}

func asCustomError(err error, ce *errors.CustomError) bool {
	if e, ok := err.(errors.CustomError); ok {
		*ce = e
		return true
	}
	return false
}

func TestRegisterHandler_RepositoryError(t *testing.T) {
	h := NewRegisterHandler(&mockUserRepo{err: errors.Internal("db error")}, &mockHasher{}, &noopTransactor{})
	err := h.Handle(context.Background(), Register{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	})
	if err == nil {
		t.Fatal("expected error for repository failure, got nil")
	}
}
