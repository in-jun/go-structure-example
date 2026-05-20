package command

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/user/domain"
	"github.com/in-jun/go-structure-example/internal/user/domain/entity"
)

const testUUID = "550e8400-e29b-41d4-a716-446655440000"

var _ domain.UserRepository = (*mockUserRepo)(nil)
var _ domain.PasswordHasher = (*mockHasher)(nil)

type mockUserRepo struct {
	user *entity.User
	err  error
}

func (m *mockUserRepo) FindByID(_ context.Context, _ string) (*entity.User, error) {
	return m.user, m.err
}
func (m *mockUserRepo) Update(_ context.Context, _ *entity.User) error { return m.err }
func (m *mockUserRepo) Delete(_ context.Context, _ string) error       { return m.err }

type mockHasher struct{}

func (m *mockHasher) Hash(password string) (string, error) { return "hashed_" + password, nil }
func (m *mockHasher) Compare(hashed, plain string) bool    { return hashed == "hashed_"+plain }

func makeUser() *entity.User {
	u, _ := entity.ReconstructUser(testUUID, "test@example.com", "hashed_oldpass", "Original", time.Now(), time.Now())
	return u
}

func TestUpdateProfileHandler_Success(t *testing.T) {
	h := NewUpdateProfileHandler(&mockUserRepo{user: makeUser()})

	err := h.Handle(context.Background(), UpdateProfile{UserID: testUUID, Name: "New Name"})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestUpdateProfileHandler_EmptyName(t *testing.T) {
	h := NewUpdateProfileHandler(&mockUserRepo{user: makeUser()})

	err := h.Handle(context.Background(), UpdateProfile{UserID: testUUID, Name: ""})
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestUpdateProfileHandler_NotFound(t *testing.T) {
	h := NewUpdateProfileHandler(&mockUserRepo{err: errors.NotFound("user not found")})

	err := h.Handle(context.Background(), UpdateProfile{UserID: testUUID, Name: "Name"})
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}
