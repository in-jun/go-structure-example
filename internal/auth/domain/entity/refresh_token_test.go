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
		wantErr   bool
	}{
		{"valid", 1, time.Now().Add(time.Hour), false},
		{"zero userID", 0, time.Now().Add(time.Hour), true},
		{"zero expiresAt", 1, time.Time{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewRefreshToken(tt.userID, tt.expiresAt)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if token.Token() == "" {
					t.Error("expected non-empty token")
				}
				if token.UserID() != tt.userID {
					t.Errorf("expected userID %d, got %d", tt.userID, token.UserID())
				}
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

func TestReconstructRefreshToken_Validation(t *testing.T) {
	now := time.Now().Add(time.Hour)

	if _, err := ReconstructRefreshToken("", 1, now); err == nil {
		t.Error("expected error for empty token")
	}
	if _, err := ReconstructRefreshToken("tok", 0, now); err == nil {
		t.Error("expected error for zero userID")
	}
	if _, err := ReconstructRefreshToken("tok", 1, time.Time{}); err == nil {
		t.Error("expected error for zero expiresAt")
	}
	rt, err := ReconstructRefreshToken("tok", 1, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rt.Token() != "tok" {
		t.Errorf("expected token 'tok', got %q", rt.Token())
	}
}
