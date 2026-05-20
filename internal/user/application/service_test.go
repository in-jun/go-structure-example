package application

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/user/application/command"
	"github.com/in-jun/go-structure-example/internal/user/application/query"
	"github.com/in-jun/go-structure-example/internal/user/domain"
	"github.com/in-jun/go-structure-example/internal/user/domain/entity"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
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

type mockHasher struct{}

func (m *mockHasher) Hash(password string) (string, error) { return "hashed_" + password, nil }
func (m *mockHasher) Compare(hashed, plain string) bool    { return hashed == "hashed_"+plain }

func newTestService(repo *mockUserRepo) *service {
	hasher := &mockHasher{}
	return NewService(
		command.NewUpdateProfileHandler(repo),
		command.NewUpdatePasswordHandler(repo, hasher),
		command.NewDeleteHandler(repo),
		query.NewGetHandler(repo),
	)
}

func makeUser() *entity.User {
	u, _ := entity.ReconstructUser(1, "test@example.com", "hashed_oldpass", "Test", time.Now(), time.Now())
	return u
}

func TestUserService_GetProfile(t *testing.T) {
	svc := newTestService(&mockUserRepo{user: makeUser()})

	result, err := svc.GetProfile(context.Background(), query.Get{UserID: 1})
	if err != nil {
		t.Fatalf("GetProfile() error = %v", err)
	}
	if result.Email != "test@example.com" {
		t.Errorf("expected test@example.com, got %q", result.Email)
	}
}

func TestUserService_GetProfile_NotFound(t *testing.T) {
	svc := newTestService(&mockUserRepo{err: errors.NotFound("user not found")})

	_, err := svc.GetProfile(context.Background(), query.Get{UserID: 99})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUserService_UpdateProfile(t *testing.T) {
	svc := newTestService(&mockUserRepo{user: makeUser()})

	err := svc.UpdateProfile(context.Background(), command.UpdateProfile{UserID: 1, Name: "New Name"})
	if err != nil {
		t.Fatalf("UpdateProfile() error = %v", err)
	}
}

func TestUserService_UpdatePassword(t *testing.T) {
	svc := newTestService(&mockUserRepo{user: makeUser()})

	err := svc.UpdatePassword(context.Background(), command.UpdatePassword{
		UserID:          1,
		CurrentPassword: "oldpass",
		NewPassword:     "newpassword123",
	})
	if err != nil {
		t.Fatalf("UpdatePassword() error = %v", err)
	}
}

func TestUserService_UpdatePassword_WrongCurrent(t *testing.T) {
	svc := newTestService(&mockUserRepo{user: makeUser()})

	err := svc.UpdatePassword(context.Background(), command.UpdatePassword{
		UserID:          1,
		CurrentPassword: "wrongpassword",
		NewPassword:     "newpassword123",
	})
	if err == nil {
		t.Fatal("expected error for wrong current password, got nil")
	}
}

func TestUserService_Delete(t *testing.T) {
	svc := newTestService(&mockUserRepo{user: makeUser()})

	err := svc.Delete(context.Background(), command.Delete{UserID: 1})
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

var _ domain.UserRepository = (*mockUserRepo)(nil)
var _ domain.PasswordHasher = (*mockHasher)(nil)
