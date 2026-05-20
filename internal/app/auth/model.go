package auth

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

type RefreshToken struct {
	token     string
	userID    uint
	expiresAt time.Time
}

func NewRefreshToken(userID uint, expiresAt time.Time) (*RefreshToken, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return &RefreshToken{
		token:     base64.URLEncoding.EncodeToString(b),
		userID:    userID,
		expiresAt: expiresAt,
	}, nil
}

func ReconstructRefreshToken(token string, userID uint, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{token: token, userID: userID, expiresAt: expiresAt}
}

func (rt *RefreshToken) Token() string        { return rt.token }
func (rt *RefreshToken) UserID() uint         { return rt.userID }
func (rt *RefreshToken) ExpiresAt() time.Time { return rt.expiresAt }
