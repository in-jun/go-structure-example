package entity

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"
)

var (
	errInvalidRefreshToken            = errors.New("user ID and expiration time are required")
	errInvalidReconstructRefreshToken = errors.New("token, user ID, and expiration time are required")
)

type RefreshToken struct {
	token     string
	userID    string
	expiresAt time.Time
}

func NewRefreshToken(userID string, expiresAt time.Time) (*RefreshToken, error) {
	if userID == "" || expiresAt.IsZero() {
		return nil, errInvalidRefreshToken
	}
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

func ReconstructRefreshToken(token, userID string, expiresAt time.Time) (*RefreshToken, error) {
	if token == "" || userID == "" || expiresAt.IsZero() {
		return nil, errInvalidReconstructRefreshToken
	}
	return &RefreshToken{token: token, userID: userID, expiresAt: expiresAt}, nil
}

func (t *RefreshToken) Token() string        { return t.token }
func (t *RefreshToken) UserID() string       { return t.userID }
func (t *RefreshToken) ExpiresAt() time.Time { return t.expiresAt }
func (t *RefreshToken) IsExpired() bool      { return time.Now().After(t.expiresAt) }
