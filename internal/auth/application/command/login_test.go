package command

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/auth/domain"
	"github.com/in-jun/go-structure-example/internal/auth/domain/entity"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

type mockTokenRepo struct {
	token *entity.RefreshToken
	userID uint
	err    error
}

func (m *mockTokenRepo) Save(_ context.Context, _ *entity.RefreshToken) error { return m.err }
func (m *mockTokenRepo) FindByToken(_ context.Context, _ string) (*entity.RefreshToken, error) {
	return m.token, m.err
}
func (m *mockTokenRepo) DeleteByUserID(_ context.Context, _ uint) error       { return m.err }
func (m *mockTokenRepo) DeleteByToken(_ context.Context, _ string) error      { return m.err }
func (m *mockTokenRepo) MarkTokenUsed(_ context.Context, _ string, _ uint) error { return m.err }
func (m *mockTokenRepo) FindUsedToken(_ context.Context, _ string) (uint, error) {
	return m.userID, m.err
}
func (m *mockTokenRepo) BlacklistAccessToken(_ context.Context, _ string, _ time.Duration) error {
	return m.err
}
func (m *mockTokenRepo) IsAccessTokenBlacklisted(_ context.Context, _ string) (bool, error) {
	return false, m.err
}
func (m *mockTokenRepo) RevokeAllAccessTokens(_ context.Context, _ uint, _ time.Duration) error {
	return m.err
}
func (m *mockTokenRepo) IsRevokedForUser(_ context.Context, _ uint, _ int64) (bool, error) {
	return false, m.err
}

type mockTokenGen struct{}

func (m *mockTokenGen) GenerateAccessToken(_ uint) (string, error) { return "access-token", nil }
func (m *mockTokenGen) ValidateToken(_ string) (*domain.TokenClaims, error) {
	return &domain.TokenClaims{UserID: 1}, nil
}
func (m *mockTokenGen) AccessExpirySeconds() int           { return 900 }
func (m *mockTokenGen) RefreshExpiry() time.Duration       { return 7 * 24 * time.Hour }

var _ domain.TokenRepository = (*mockTokenRepo)(nil)
var _ domain.TokenGenerator = (*mockTokenGen)(nil)

func makeAuthUser() *entity.User {
	u, _ := entity.NewUser("test@example.com", "hashed_password123", "Test User")
	u.SetID(1)
	return u
}

func TestLoginHandler_Success(t *testing.T) {
	h := NewLoginHandler(&mockUserRepo{user: makeAuthUser()}, &mockTokenRepo{}, &mockTokenGen{}, &mockHasher{})

	result, err := h.Handle(context.Background(), Login{Email: "test@example.com", Password: "password123"})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if result.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if result.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
}

func TestLoginHandler_UserNotFound(t *testing.T) {
	h := NewLoginHandler(&mockUserRepo{}, &mockTokenRepo{}, &mockTokenGen{}, &mockHasher{})

	_, err := h.Handle(context.Background(), Login{Email: "notfound@example.com", Password: "password123"})
	if err == nil {
		t.Fatal("expected error for unknown user, got nil")
	}
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	h := NewLoginHandler(&mockUserRepo{user: makeAuthUser()}, &mockTokenRepo{}, &mockTokenGen{}, &mockHasher{})

	_, err := h.Handle(context.Background(), Login{Email: "test@example.com", Password: "wrongpassword"})
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
	var ce errors.CustomError
	if !asCustomError(err, &ce) || ce.Status != 401 {
		t.Errorf("expected 401 Unauthorized, got %v", err)
	}
}

func TestLoginHandler_InvalidEmail(t *testing.T) {
	h := NewLoginHandler(&mockUserRepo{}, &mockTokenRepo{}, &mockTokenGen{}, &mockHasher{})

	_, err := h.Handle(context.Background(), Login{Email: "not-an-email", Password: "password123"})
	if err == nil {
		t.Fatal("expected error for invalid email, got nil")
	}
}

func TestLoginHandler_RepositoryError(t *testing.T) {
	h := NewLoginHandler(&mockUserRepo{err: errors.Internal("db error")}, &mockTokenRepo{}, &mockTokenGen{}, &mockHasher{})

	_, err := h.Handle(context.Background(), Login{Email: "test@example.com", Password: "password123"})
	if err == nil {
		t.Fatal("expected error for repository failure, got nil")
	}
}
