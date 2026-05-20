package entity

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		password  string
		userName  string
		wantError bool
	}{
		{"valid", "test@example.com", "hashed_password", "Test User", false},
		{"empty email", "", "hashed_password", "Test User", true},
		{"empty password", "test@example.com", "", "Test User", true},
		{"empty name", "test@example.com", "hashed_password", "", true},
		{"all empty", "", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.email, tt.password, tt.userName)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", user)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if !tt.wantError && user.ID() != 0 {
				t.Error("expected zero ID for new user before save")
			}
		})
	}
}

func TestReconstructUser(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		id        uint
		email     string
		password  string
		userName  string
		wantError bool
	}{
		{"valid", 1, "test@example.com", "hashed", "Test", false},
		{"zero id", 0, "test@example.com", "hashed", "Test", true},
		{"empty email", 1, "", "hashed", "Test", true},
		{"empty password", 1, "test@example.com", "", "Test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := ReconstructUser(tt.id, tt.email, tt.password, tt.userName, now, now)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", user)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestUser_SetID(t *testing.T) {
	u, _ := NewUser("test@example.com", "hashed", "Test")
	if u.ID() != 0 {
		t.Error("expected zero ID before set")
	}
	u.SetID(42)
	if u.ID() != 42 {
		t.Errorf("expected ID 42, got %d", u.ID())
	}
}
