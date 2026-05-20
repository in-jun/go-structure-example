package auth

import (
	"context"
	"time"
)

type Repository interface {
	Save(ctx context.Context, token *RefreshToken) error
	FindByToken(ctx context.Context, token string) (*RefreshToken, error)
	DeleteByUserID(ctx context.Context, userID uint) error
	DeleteByToken(ctx context.Context, token string) error
	BlacklistAccessToken(ctx context.Context, jti string, ttl time.Duration) error
	IsAccessTokenBlacklisted(ctx context.Context, jti string) (bool, error)
}
