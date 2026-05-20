package entity

import (
	"testing"
	"time"
)

const testUUID = "550e8400-e29b-41d4-a716-446655440000"

func TestReconstructUser(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		id        string
		email     string
		userName  string
		wantError bool
	}{
		{"valid", testUUID, "test@example.com", "Test User", false},
		{"empty id", "", "test@example.com", "Test User", true},
		{"empty email", testUUID, "", "Test User", true},
		{"empty name", testUUID, "test@example.com", "", true},
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
	u, _ := ReconstructUser(testUUID, "test@example.com", "hashed", "Original", now, now)

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
	u, _ := ReconstructUser(testUUID, "test@example.com", "old_hash", "Test", now, now)

	u.SetPassword("new_hash")
	if u.HashedPassword() != "new_hash" {
		t.Errorf("expected new_hash, got %q", u.HashedPassword())
	}
}
