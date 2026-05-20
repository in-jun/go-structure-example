package domain

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/auth/domain/entity"
)

type UserRepository interface {
	Save(ctx context.Context, user *entity.User) error
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

type TokenRepository interface {
	Save(ctx context.Context, token *entity.RefreshToken) error
	FindByToken(ctx context.Context, token string) (*entity.RefreshToken, error)
	DeleteByUserID(ctx context.Context, userID uint) error
	DeleteByToken(ctx context.Context, token string) error
	BlacklistAccessToken(ctx context.Context, jti string, ttl time.Duration) error
	IsAccessTokenBlacklisted(ctx context.Context, jti string) (bool, error)
	RevokeAllAccessTokens(ctx context.Context, userID uint, ttl time.Duration) error
	IsRevokedForUser(ctx context.Context, userID uint, issuedAt int64) (bool, error)
}

type TokenClaims struct {
	UserID   uint
	JTI      string
	IssuedAt int64
}

type TokenGenerator interface {
	GenerateAccessToken(userID uint) (string, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
	AccessExpirySeconds() int
	RefreshExpiry() time.Duration
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) bool
}
