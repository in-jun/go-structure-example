package jwt

import (
	"testing"
	"time"
)

func newTestProvider(t *testing.T) *provider {
	t.Helper()
	p, err := NewProvider("test-secret-key", "15m", "168h")
	if err != nil {
		t.Fatalf("NewProvider() error = %v", err)
	}
	return p.(*provider)
}

func TestProvider_GenerateAndValidate(t *testing.T) {
	p := newTestProvider(t)
	const userID = uint(42)

	token, err := p.GenerateAccessToken(userID)
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := p.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("expected userID %d, got %d", userID, claims.UserID)
	}
	if claims.JTI == "" {
		t.Error("expected non-empty JTI")
	}
	if claims.IssuedAt == 0 {
		t.Error("expected non-zero IssuedAt")
	}
}

func TestProvider_ValidateToken_InvalidSignature(t *testing.T) {
	p := newTestProvider(t)

	p2, _ := NewProvider("different-secret", "15m", "168h")
	token, _ := p2.(*provider).GenerateAccessToken(1)

	_, err := p.ValidateToken(token)
	if err == nil {
		t.Error("expected error for token with wrong signature")
	}
}

func TestProvider_ValidateToken_Malformed(t *testing.T) {
	p := newTestProvider(t)

	_, err := p.ValidateToken("not.a.valid.jwt")
	if err == nil {
		t.Error("expected error for malformed token")
	}
}

func TestProvider_UniqueJTI(t *testing.T) {
	p := newTestProvider(t)

	c1, _ := p.ValidateToken(func() string { t, _ := p.GenerateAccessToken(1); return t }())
	c2, _ := p.ValidateToken(func() string { t, _ := p.GenerateAccessToken(1); return t }())

	if c1.JTI == c2.JTI {
		t.Error("expected unique JTI per token")
	}
}

func TestProvider_AccessExpirySeconds(t *testing.T) {
	p := newTestProvider(t)
	if p.AccessExpirySeconds() != int((15 * time.Minute).Seconds()) {
		t.Errorf("expected %d, got %d", int((15*time.Minute).Seconds()), p.AccessExpirySeconds())
	}
}

func TestProvider_RefreshExpiry(t *testing.T) {
	p := newTestProvider(t)
	if p.RefreshExpiry() != 168*time.Hour {
		t.Errorf("expected 168h, got %v", p.RefreshExpiry())
	}
}

func TestNewProvider_InvalidDuration(t *testing.T) {
	_, err := NewProvider("secret", "invalid", "168h")
	if err == nil {
		t.Error("expected error for invalid access expiry")
	}

	_, err = NewProvider("secret", "15m", "notaduration")
	if err == nil {
		t.Error("expected error for invalid refresh expiry")
	}
}
