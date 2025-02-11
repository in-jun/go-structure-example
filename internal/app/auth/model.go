package auth

import (
	"time"
)

type RefreshToken struct {
	RefreshToken string    `json:"refresh_token"`
	UserID       uint      `json:"user_id"`
	ExpiresAt    time.Time `json:"expires_at"`
}
