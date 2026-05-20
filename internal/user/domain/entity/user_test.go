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
		uname     string
		wantErr   bool
	}{
		{"valid", 1, "test@example.com", "Test User", false},
		{"zero id", 0, "test@example.com", "Test User", true},
		{"empty email", 1, "", "Test User", true},
		{"empty name", 1, "test@example.com", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := ReconstructUser(tt.id, tt.email, "hashed", tt.uname, now, now)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if u.ID() != tt.id {
					t.Errorf("expected ID %d, got %d", tt.id, u.ID())
				}
				if u.Email() != tt.email {
					t.Errorf("expected email %q, got %q", tt.email, u.Email())
				}
				if u.Name() != tt.uname {
					t.Errorf("expected name %q, got %q", tt.uname, u.Name())
				}
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
