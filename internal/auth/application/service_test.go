package application

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/auth/application/command"
	"github.com/in-jun/go-structure-example/internal/auth/application/query"
	"github.com/in-jun/go-structure-example/internal/auth/domain"
	"github.com/in-jun/go-structure-example/internal/auth/domain/entity"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

type noopTransactor struct{}

func (n *noopTransactor) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

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

type mockTokenRepo struct {
	token *entity.RefreshToken
	err   error
}

func (m *mockTokenRepo) Save(_ context.Context, _ *entity.RefreshToken) error       { return m.err }
func (m *mockTokenRepo) FindByToken(_ context.Context, _ string) (*entity.RefreshToken, error) {
	return m.token, m.err
}
func (m *mockTokenRepo) DeleteByUserID(_ context.Context, _ uint) error  { return m.err }
func (m *mockTokenRepo) DeleteByToken(_ context.Context, _ string) error { return m.err }
func (m *mockTokenRepo) MarkTokenUsed(_ context.Context, _ string, _ uint) error { return m.err }
func (m *mockTokenRepo) FindUsedToken(_ context.Context, _ string) (uint, error) {
	return 0, m.err
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

var _ domain.TokenRepository = (*mockTokenRepo)(nil)
var _ domain.UserRepository = (*mockUserRepo)(nil)

type mockTokenGen struct {
	accessToken string
	claims      *domain.TokenClaims
	err         error
}

func (m *mockTokenGen) GenerateAccessToken(_ uint) (string, error) {
	return m.accessToken, m.err
}
func (m *mockTokenGen) ValidateToken(_ string) (*domain.TokenClaims, error) {
	return m.claims, m.err
}
func (m *mockTokenGen) AccessExpirySeconds() int    { return 3600 }
func (m *mockTokenGen) RefreshExpiry() time.Duration { return 24 * time.Hour }

type mockHasher struct{}

func (m *mockHasher) Hash(password string) (string, error) { return "hashed_" + password, nil }
func (m *mockHasher) Compare(hashed, plain string) bool    { return hashed == "hashed_"+plain }

func newTestService(userRepo *mockUserRepo, tokenRepo *mockTokenRepo, tokenGen *mockTokenGen) *service {
	hasher := &mockHasher{}
	return NewService(
		command.NewRegisterHandler(userRepo, hasher, &noopTransactor{}),
		command.NewLoginHandler(userRepo, tokenRepo, tokenGen, hasher),
		command.NewRefreshHandler(tokenRepo, tokenGen),
		command.NewLogoutHandler(tokenRepo, tokenGen),
		command.NewLogoutAllHandler(tokenRepo, tokenGen),
		query.NewValidateHandler(tokenRepo, tokenGen),
	)
}

func TestAuthService_Register(t *testing.T) {
	svc := newTestService(&mockUserRepo{}, &mockTokenRepo{}, &mockTokenGen{})

	err := svc.Register(context.Background(), command.Register{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
	})
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	existingUser, _ := entity.ReconstructUser(1, "test@example.com", "hashed", "Existing", time.Now(), time.Now())
	svc := newTestService(&mockUserRepo{user: existingUser}, &mockTokenRepo{}, &mockTokenGen{})

	err := svc.Register(context.Background(), command.Register{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test",
	})
	if err == nil {
		t.Fatal("expected error for duplicate email, got nil")
	}
}

func TestAuthService_Login(t *testing.T) {
	user, _ := entity.ReconstructUser(1, "test@example.com", "hashed_password123", "Test", time.Now(), time.Now())
	tokenGen := &mockTokenGen{accessToken: "access-token"}
	svc := newTestService(&mockUserRepo{user: user}, &mockTokenRepo{}, tokenGen)

	result, err := svc.Login(context.Background(), command.Login{
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if result.AccessToken != "access-token" {
		t.Errorf("expected access-token, got %q", result.AccessToken)
	}
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	user, _ := entity.ReconstructUser(1, "test@example.com", "hashed_correct", "Test", time.Now(), time.Now())
	svc := newTestService(&mockUserRepo{user: user}, &mockTokenRepo{}, &mockTokenGen{})

	_, err := svc.Login(context.Background(), command.Login{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	svc := newTestService(&mockUserRepo{err: errors.NotFound("user not found")}, &mockTokenRepo{}, &mockTokenGen{})

	_, err := svc.Login(context.Background(), command.Login{
		Email:    "notfound@example.com",
		Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for user not found, got nil")
	}
}

type reuseDetectTokenRepo struct {
	mockTokenRepo
}

func (m *reuseDetectTokenRepo) FindUsedToken(_ context.Context, _ string) (uint, error) {
	return 42, nil // non-zero: token was previously used (theft attempt)
}

func newServiceWithRepo(tokenRepo domain.TokenRepository) *service {
	tokenGen := &mockTokenGen{}
	return &service{
		register:  command.NewRegisterHandler(&mockUserRepo{}, &mockHasher{}, &noopTransactor{}),
		login:     command.NewLoginHandler(&mockUserRepo{}, tokenRepo, tokenGen, &mockHasher{}),
		refresh:   command.NewRefreshHandler(tokenRepo, tokenGen),
		logout:    command.NewLogoutHandler(tokenRepo, tokenGen),
		logoutAll: command.NewLogoutAllHandler(tokenRepo, tokenGen),
		validate:  query.NewValidateHandler(tokenRepo, tokenGen),
	}
}

func TestAuthService_Refresh_TokenReuseDetected(t *testing.T) {
	svc := newServiceWithRepo(&reuseDetectTokenRepo{})

	_, err := svc.Refresh(context.Background(), command.Refresh{RefreshToken: "already-used-token"})
	if err == nil {
		t.Fatal("expected error for token reuse, got nil")
	}
}

func TestAuthService_Logout(t *testing.T) {
	tokenRepo := &mockTokenRepo{}
	tokenGen := &mockTokenGen{claims: &domain.TokenClaims{UserID: 1, JTI: "jti"}}
	svc := newTestService(&mockUserRepo{}, tokenRepo, tokenGen)

	err := svc.Logout(context.Background(), command.Logout{
		RefreshToken:   "refresh-token",
		AccessTokenJTI: "jti",
	})
	if err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
}

func TestAuthService_LogoutAll(t *testing.T) {
	svc := newTestService(&mockUserRepo{}, &mockTokenRepo{}, &mockTokenGen{})

	err := svc.LogoutAll(context.Background(), command.LogoutAll{UserID: 1})
	if err != nil {
		t.Fatalf("LogoutAll() error = %v", err)
	}
}
