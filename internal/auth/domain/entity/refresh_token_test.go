package entity

import (
	"testing"
	"time"
)

func TestNewRefreshToken(t *testing.T) {
	tests := []struct {
		name      string
		userID    uint
		expiresAt time.Time
		wantError bool
	}{
		{"valid", 1, time.Now().Add(time.Hour), false},
		{"zero userID", 0, time.Now().Add(time.Hour), true},
		{"zero expiresAt", 1, time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewRefreshToken(tt.userID, tt.expiresAt)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", token)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if !tt.wantError && token.Token() == "" {
				t.Error("expected non-empty token")
			}
			if !tt.wantError && token.UserID() != tt.userID {
				t.Errorf("expected userID %d, got %d", tt.userID, token.UserID())
			}
		})
	}
}

func TestRefreshToken_IsExpired(t *testing.T) {
	expired, _ := ReconstructRefreshToken("tok", 1, time.Now().Add(-time.Hour))
	if !expired.IsExpired() {
		t.Error("expected token to be expired")
	}

	valid, _ := ReconstructRefreshToken("tok", 1, time.Now().Add(time.Hour))
	if valid.IsExpired() {
		t.Error("expected token to be valid")
	}
}

func TestReconstructRefreshToken(t *testing.T) {
	now := time.Now().Add(time.Hour)

	tests := []struct {
		name      string
		token     string
		userID    uint
		expiresAt time.Time
		wantError bool
	}{
		{"valid", "tok", 1, now, false},
		{"empty token", "", 1, now, true},
		{"zero userID", "tok", 0, now, true},
		{"zero expiresAt", "tok", 1, time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt, err := ReconstructRefreshToken(tt.token, tt.userID, tt.expiresAt)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", rt)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if !tt.wantError && rt.Token() != tt.token {
				t.Errorf("Token = %q, want %q", rt.Token(), tt.token)
			}
		})
	}
}
