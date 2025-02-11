package auth

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, refreshToken *RefreshToken) error
	FindByToken(ctx context.Context, refreshToken string) (*RefreshToken, error)
	DeleteByUserId(ctx context.Context, userID uint) error
	DeleteByToken(ctx context.Context, refreshToken string) error
}
