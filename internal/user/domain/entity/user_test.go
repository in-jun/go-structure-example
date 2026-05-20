package entity

import (
	"testing"
	"time"
)

func TestReconstructUser(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		id        uint
		email     string
		userName  string
		wantError bool
	}{
		{"valid", 1, "test@example.com", "Test User", false},
		{"zero id", 0, "test@example.com", "Test User", true},
		{"empty email", 1, "", "Test User", true},
		{"empty name", 1, "test@example.com", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := ReconstructUser(tt.id, tt.email, "hashed", tt.userName, now, now)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", user)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestUser_SetName(t *testing.T) {
	now := time.Now()
	u, _ := ReconstructUser(1, "test@example.com", "hashed", "Original", now, now)

	u.SetName("Updated")
	if u.Name() != "Updated" {
		t.Errorf("expected Updated, got %q", u.Name())
	}
	if !u.UpdatedAt().After(now) {
		t.Error("expected updatedAt to be updated")
	}
}

func TestUser_SetPassword(t *testing.T) {
	now := time.Now()
	u, _ := ReconstructUser(1, "test@example.com", "old_hash", "Test", now, now)

	u.SetPassword("new_hash")
	if u.HashedPassword() != "new_hash" {
		t.Errorf("expected new_hash, got %q", u.HashedPassword())
	}
}
