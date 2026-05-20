package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/in-jun/go-structure-example/internal/app/auth"
	appErrors "github.com/in-jun/go-structure-example/internal/pkg/utils/errors"

	goredis "github.com/go-redis/redis/v8"
)

type authRepository struct {
	client *goredis.Client
}

func NewAuthRepository(client *goredis.Client) auth.Repository {
	return &authRepository{client: client}
}

func (r *authRepository) tokenKey(token string) string {
	return fmt.Sprintf("rt:lookup:%s", token)
}

func (r *authRepository) userKey(userID uint) string {
	return fmt.Sprintf("rt:user:%d", userID)
}

func (r *authRepository) blacklistKey(jti string) string {
	return fmt.Sprintf("bl:jti:%s", jti)
}

func (r *authRepository) Save(ctx context.Context, token *auth.RefreshToken) error {
	raw := struct {
		Token     string    `json:"token"`
		UserID    uint      `json:"user_id"`
		ExpiresAt time.Time `json:"expires_at"`
	}{
		Token:     token.Token(),
		UserID:    token.UserID(),
		ExpiresAt: token.ExpiresAt(),
	}

	data, err := json.Marshal(raw)
	if err != nil {
		return appErrors.Internal("Failed to marshal refresh token")
	}

	expiration := time.Until(token.ExpiresAt())

	pipe := r.client.Pipeline()
	pipe.Set(ctx, r.tokenKey(token.Token()), data, expiration)
	pipe.SAdd(ctx, r.userKey(token.UserID()), token.Token())
	pipe.Expire(ctx, r.userKey(token.UserID()), expiration)

	if _, err := pipe.Exec(ctx); err != nil {
		return appErrors.Internal("Failed to save refresh token")
	}

	return nil
}

func (r *authRepository) FindByToken(ctx context.Context, tokenStr string) (*auth.RefreshToken, error) {
	data, err := r.client.Get(ctx, r.tokenKey(tokenStr)).Result()
	if errors.Is(err, goredis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, appErrors.Internal("Failed to get refresh token")
	}

	var raw struct {
		Token     string    `json:"token"`
		UserID    uint      `json:"user_id"`
		ExpiresAt time.Time `json:"expires_at"`
	}
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return nil, appErrors.Internal("Failed to unmarshal refresh token")
	}

	return auth.ReconstructRefreshToken(raw.Token, raw.UserID, raw.ExpiresAt), nil
}

func (r *authRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	tokens, err := r.client.SMembers(ctx, r.userKey(userID)).Result()
	if errors.Is(err, goredis.Nil) || len(tokens) == 0 {
		return nil
	}
	if err != nil {
		return appErrors.Internal("Failed to get user tokens")
	}

	pipe := r.client.Pipeline()
	for _, token := range tokens {
		pipe.Del(ctx, r.tokenKey(token))
	}
	pipe.Del(ctx, r.userKey(userID))

	if _, err := pipe.Exec(ctx); err != nil {
		return appErrors.Internal("Failed to delete refresh tokens")
	}

	return nil
}

func (r *authRepository) DeleteByToken(ctx context.Context, tokenStr string) error {
	data, err := r.client.Get(ctx, r.tokenKey(tokenStr)).Result()
	if errors.Is(err, goredis.Nil) {
		return nil
	}
	if err != nil {
		return appErrors.Internal("Failed to get refresh token for deletion")
	}

	var raw struct {
		UserID uint `json:"user_id"`
	}
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return appErrors.Internal("Failed to unmarshal refresh token")
	}

	pipe := r.client.Pipeline()
	pipe.SRem(ctx, r.userKey(raw.UserID), tokenStr)
	pipe.Del(ctx, r.tokenKey(tokenStr))

	if _, err := pipe.Exec(ctx); err != nil {
		return appErrors.Internal("Failed to delete refresh token")
	}

	return nil
}

func (r *authRepository) BlacklistAccessToken(ctx context.Context, jti string, ttl time.Duration) error {
	if err := r.client.Set(ctx, r.blacklistKey(jti), "1", ttl).Err(); err != nil {
		return appErrors.Internal("Failed to blacklist access token")
	}
	return nil
}

func (r *authRepository) IsAccessTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	_, err := r.client.Get(ctx, r.blacklistKey(jti)).Result()
	if errors.Is(err, goredis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, appErrors.Internal("Failed to check token blacklist")
	}
	return true, nil
}
