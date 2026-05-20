package query

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/auth/domain"
	"github.com/in-jun/go-structure-example/internal/auth/domain/entity"
)

type mockTokenRepo struct {
	blacklisted bool
	revoked     bool
	err         error
}

func (m *mockTokenRepo) Save(_ context.Context, _ *entity.RefreshToken) error { return m.err }
func (m *mockTokenRepo) FindByToken(_ context.Context, _ string) (*entity.RefreshToken, error) {
	return nil, m.err
}
func (m *mockTokenRepo) DeleteByUserID(_ context.Context, _ uint) error            { return m.err }
func (m *mockTokenRepo) DeleteByToken(_ context.Context, _ string) error           { return m.err }
func (m *mockTokenRepo) MarkTokenUsed(_ context.Context, _ string, _ uint) error   { return m.err }
func (m *mockTokenRepo) FindUsedToken(_ context.Context, _ string) (uint, error)   { return 0, m.err }
func (m *mockTokenRepo) BlacklistAccessToken(_ context.Context, _ string, _ time.Duration) error {
	return m.err
}
func (m *mockTokenRepo) IsAccessTokenBlacklisted(_ context.Context, _ string) (bool, error) {
	return m.blacklisted, m.err
}
func (m *mockTokenRepo) RevokeAllAccessTokens(_ context.Context, _ uint, _ time.Duration) error {
	return m.err
}
func (m *mockTokenRepo) IsRevokedForUser(_ context.Context, _ uint, _ int64) (bool, error) {
	return m.revoked, m.err
}

var _ domain.TokenRepository = (*mockTokenRepo)(nil)
var _ domain.TokenGenerator = (*mockTokenGen)(nil)

type mockTokenGen struct {
	claims *domain.TokenClaims
	err    error
}

func (m *mockTokenGen) GenerateAccessToken(_ uint) (string, error) { return "", nil }
func (m *mockTokenGen) ValidateToken(_ string) (*domain.TokenClaims, error) {
	return m.claims, m.err
}
func (m *mockTokenGen) AccessExpirySeconds() int      { return 3600 }
func (m *mockTokenGen) RefreshExpiry() time.Duration  { return 24 * time.Hour }

func TestValidateHandler_ValidToken(t *testing.T) {
	claims := &domain.TokenClaims{UserID: 1, JTI: "jti-1", IssuedAt: time.Now().Unix()}
	h := NewValidateHandler(&mockTokenRepo{}, &mockTokenGen{claims: claims})

	result, err := h.Handle(context.Background(), Validate{TokenString: "valid-token"})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if result.UserID != 1 {
		t.Errorf("expected UserID 1, got %d", result.UserID)
	}
	if result.JTI != "jti-1" {
		t.Errorf("expected JTI 'jti-1', got %q", result.JTI)
	}
}

func TestValidateHandler_EmptyToken(t *testing.T) {
	h := NewValidateHandler(&mockTokenRepo{}, &mockTokenGen{})

	_, err := h.Handle(context.Background(), Validate{TokenString: ""})
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestValidateHandler_InvalidToken(t *testing.T) {
	h := NewValidateHandler(&mockTokenRepo{}, &mockTokenGen{err: stderrors.New("token invalid")})

	_, err := h.Handle(context.Background(), Validate{TokenString: "bad-token"})
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
}

func TestValidateHandler_BlacklistedToken(t *testing.T) {
	claims := &domain.TokenClaims{UserID: 1, JTI: "jti-bl", IssuedAt: time.Now().Unix()}
	h := NewValidateHandler(&mockTokenRepo{blacklisted: true}, &mockTokenGen{claims: claims})

	_, err := h.Handle(context.Background(), Validate{TokenString: "blacklisted-token"})
	if err == nil {
		t.Fatal("expected error for blacklisted token, got nil")
	}
}

func TestValidateHandler_RevokedForUser(t *testing.T) {
	claims := &domain.TokenClaims{UserID: 2, JTI: "jti-rv", IssuedAt: time.Now().Unix()}
	h := NewValidateHandler(&mockTokenRepo{revoked: true}, &mockTokenGen{claims: claims})

	_, err := h.Handle(context.Background(), Validate{TokenString: "revoked-token"})
	if err == nil {
		t.Fatal("expected error for revoked token, got nil")
	}
}
