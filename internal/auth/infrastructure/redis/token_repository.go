package redis

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/in-jun/go-structure-example/internal/auth/domain"
	"github.com/in-jun/go-structure-example/internal/auth/domain/entity"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

var _ domain.TokenRepository = (*tokenRepository)(nil)

type tokenRepository struct {
	client *goredis.Client
}

func NewTokenRepository(client *goredis.Client) domain.TokenRepository {
	return &tokenRepository{client: client}
}

func (r *tokenRepository) tokenKey(token string) string {
	return fmt.Sprintf("rt:lookup:%s", token)
}

func (r *tokenRepository) userKey(userID uint) string {
	return fmt.Sprintf("rt:user:%d", userID)
}

func (r *tokenRepository) blacklistKey(jti string) string {
	return fmt.Sprintf("bl:jti:%s", jti)
}

func (r *tokenRepository) Save(ctx context.Context, token *entity.RefreshToken) error {
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
		return errors.Internal("Failed to marshal refresh token")
	}

	expiration := time.Until(token.ExpiresAt())
	pipe := r.client.Pipeline()
	pipe.Set(ctx, r.tokenKey(token.Token()), data, expiration)
	pipe.SAdd(ctx, r.userKey(token.UserID()), token.Token())
	pipe.Expire(ctx, r.userKey(token.UserID()), expiration)

	if _, err := pipe.Exec(ctx); err != nil {
		return errors.Internal("Failed to save refresh token")
	}
	return nil
}

func (r *tokenRepository) FindByToken(ctx context.Context, tokenStr string) (*entity.RefreshToken, error) {
	data, err := r.client.Get(ctx, r.tokenKey(tokenStr)).Result()
	if stderrors.Is(err, goredis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get refresh token")
	}

	var raw struct {
		Token     string    `json:"token"`
		UserID    uint      `json:"user_id"`
		ExpiresAt time.Time `json:"expires_at"`
	}
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return nil, errors.Internal("Failed to unmarshal refresh token")
	}
	rt, err := entity.ReconstructRefreshToken(raw.Token, raw.UserID, raw.ExpiresAt)
	if err != nil {
		return nil, errors.Internal("Failed to reconstruct refresh token")
	}
	return rt, nil
}

func (r *tokenRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	tokens, err := r.client.SMembers(ctx, r.userKey(userID)).Result()
	if stderrors.Is(err, goredis.Nil) || len(tokens) == 0 {
		return nil
	}
	if err != nil {
		return errors.Internal("Failed to get user tokens")
	}

	pipe := r.client.Pipeline()
	for _, token := range tokens {
		pipe.Del(ctx, r.tokenKey(token))
	}
	pipe.Del(ctx, r.userKey(userID))

	if _, err := pipe.Exec(ctx); err != nil {
		return errors.Internal("Failed to delete refresh tokens")
	}
	return nil
}

func (r *tokenRepository) DeleteByToken(ctx context.Context, tokenStr string) error {
	data, err := r.client.Get(ctx, r.tokenKey(tokenStr)).Result()
	if stderrors.Is(err, goredis.Nil) {
		return nil
	}
	if err != nil {
		return errors.Internal("Failed to get refresh token for deletion")
	}

	var raw struct {
		UserID uint `json:"user_id"`
	}
	if err := json.Unmarshal([]byte(data), &raw); err != nil {
		return errors.Internal("Failed to unmarshal refresh token")
	}

	pipe := r.client.Pipeline()
	pipe.SRem(ctx, r.userKey(raw.UserID), tokenStr)
	pipe.Del(ctx, r.tokenKey(tokenStr))

	if _, err := pipe.Exec(ctx); err != nil {
		return errors.Internal("Failed to delete refresh token")
	}
	return nil
}

const usedTokenTTL = 7 * 24 * time.Hour

func (r *tokenRepository) usedKey(token string) string {
	return fmt.Sprintf("rt:used:%s", token)
}

func (r *tokenRepository) MarkTokenUsed(ctx context.Context, token string, userID uint) error {
	if err := r.client.Set(ctx, r.usedKey(token), fmt.Sprintf("%d", userID), usedTokenTTL).Err(); err != nil {
		return errors.Internal("Failed to mark token as used")
	}
	return nil
}

func (r *tokenRepository) FindUsedToken(ctx context.Context, token string) (uint, error) {
	val, err := r.client.Get(ctx, r.usedKey(token)).Result()
	if stderrors.Is(err, goredis.Nil) {
		return 0, nil
	}
	if err != nil {
		return 0, errors.Internal("Failed to check used token")
	}
	var userID uint
	if _, err := fmt.Sscanf(val, "%d", &userID); err != nil {
		return 0, errors.Internal("Failed to parse used token user ID")
	}
	return userID, nil
}

func (r *tokenRepository) BlacklistAccessToken(ctx context.Context, jti string, ttl time.Duration) error {
	if err := r.client.Set(ctx, r.blacklistKey(jti), "1", ttl).Err(); err != nil {
		return errors.Internal("Failed to blacklist access token")
	}
	return nil
}

func (r *tokenRepository) IsAccessTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	_, err := r.client.Get(ctx, r.blacklistKey(jti)).Result()
	if stderrors.Is(err, goredis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, errors.Internal("Failed to check token blacklist")
	}
	return true, nil
}

func (r *tokenRepository) revokedAtKey(userID uint) string {
	return fmt.Sprintf("revoked_at:user:%d", userID)
}

func (r *tokenRepository) RevokeAllAccessTokens(ctx context.Context, userID uint, ttl time.Duration) error {
	now := fmt.Sprintf("%d", time.Now().Unix())
	if err := r.client.Set(ctx, r.revokedAtKey(userID), now, ttl).Err(); err != nil {
		return errors.Internal("Failed to revoke all access tokens")
	}
	return nil
}

func (r *tokenRepository) IsRevokedForUser(ctx context.Context, userID uint, issuedAt int64) (bool, error) {
	val, err := r.client.Get(ctx, r.revokedAtKey(userID)).Result()
	if stderrors.Is(err, goredis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, errors.Internal("Failed to check token revocation")
	}
	var revokedAt int64
	if _, err := fmt.Sscanf(val, "%d", &revokedAt); err != nil {
		return false, errors.Internal("Failed to parse revocation timestamp")
	}
	return issuedAt <= revokedAt, nil
}
