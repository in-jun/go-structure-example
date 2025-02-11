package redis

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/in-jun/go-structure-example/internal/app/auth"
	"github.com/in-jun/go-structure-example/internal/pkg/utils/errors"

	"github.com/go-redis/redis/v8"
)

type authRepository struct {
	redis *redis.Client
}

func NewAuthRepository(redis *redis.Client) auth.Repository {
	return &authRepository{redis: redis}
}

func (r *authRepository) Save(ctx context.Context, token *auth.RefreshToken) error {
	key := "refresh_token:" + strconv.FormatUint(uint64(token.UserID), 10)
	value, err := json.Marshal(token)
	if err != nil {
		return errors.Internal("Failed to marshal refresh token")
	}

	expiration := time.Until(token.ExpiresAt)
	if err := r.redis.Set(ctx, key, value, expiration).Err(); err != nil {
		return errors.Internal("Failed to save refresh token")
	}

	return nil
}

func (r *authRepository) FindByUserId(ctx context.Context, userID uint) (*auth.RefreshToken, error) {
	key := "refresh_token:" + strconv.FormatUint(uint64(userID), 10)
	value, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get refresh token")
	}

	var token auth.RefreshToken
	if err := json.Unmarshal([]byte(value), &token); err != nil {
		return nil, errors.Internal("Failed to unmarshal refresh token")
	}

	return &token, nil
}

func (r *authRepository) FindByToken(ctx context.Context, refreshToken string) (*auth.RefreshToken, error) {
	pattern := "refresh_token:*"
	iter := r.redis.Scan(ctx, 0, pattern, 0).Iterator()

	for iter.Next(ctx) {
		key := iter.Val()
		value, err := r.redis.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var token auth.RefreshToken
		if err := json.Unmarshal([]byte(value), &token); err != nil {
			continue
		}

		if token.RefreshToken == refreshToken {
			return &token, nil
		}
	}

	return nil, nil
}

func (r *authRepository) DeleteByUserId(ctx context.Context, userID uint) error {
	key := "refresh_token:" + strconv.FormatUint(uint64(userID), 10)
	if err := r.redis.Del(ctx, key).Err(); err != nil {
		return errors.Internal("Failed to delete refresh token")
	}
	return nil
}

func (r *authRepository) DeleteByToken(ctx context.Context, refreshToken string) error {
	token, err := r.FindByToken(ctx, refreshToken)
	if err != nil {
		return err
	}
	if token == nil {
		return errors.NotFound("Refresh token not found")
	}

	return r.DeleteByUserId(ctx, token.UserID)
}

// token blacklist 관련
func (r *authRepository) AddToBlacklist(ctx context.Context, tokenID string, expiration time.Duration) error {
	key := "blacklist:" + tokenID
	if err := r.redis.Set(ctx, key, "1", expiration).Err(); err != nil {
		return errors.Internal("Failed to add token to blacklist")
	}
	return nil
}

func (r *authRepository) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := "blacklist:" + tokenID
	_, err := r.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, errors.Internal("Failed to check blacklist")
	}
	return true, nil
}
