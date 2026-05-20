package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/in-jun/go-structure-example/internal/auth/domain"
)

var _ domain.TokenGenerator = (*provider)(nil)

type claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type provider struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewProvider(secret, accessExpiry, refreshExpiry string) (domain.TokenGenerator, error) {
	access, err := time.ParseDuration(accessExpiry)
	if err != nil {
		return nil, fmt.Errorf("invalid access expiry %q: %w", accessExpiry, err)
	}
	refresh, err := time.ParseDuration(refreshExpiry)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh expiry %q: %w", refreshExpiry, err)
	}
	return &provider{
		secretKey:     []byte(secret),
		accessExpiry:  access,
		refreshExpiry: refresh,
	}, nil
}

func (p *provider) GenerateAccessToken(userID string) (string, error) {
	jti, err := generateRandomString(16)
	if err != nil {
		return "", err
	}
	c := claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(p.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(p.secretKey)
}

func (p *provider) ValidateToken(tokenString string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return p.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := token.Claims.(*claims); ok && token.Valid {
		var issuedAt int64
		if c.IssuedAt != nil {
			issuedAt = c.IssuedAt.Unix()
		}
		return &domain.TokenClaims{UserID: c.UserID, JTI: c.ID, IssuedAt: issuedAt}, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

func (p *provider) AccessExpirySeconds() int {
	return int(p.accessExpiry.Seconds())
}

func (p *provider) RefreshExpiry() time.Duration {
	return p.refreshExpiry
}

func generateRandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
