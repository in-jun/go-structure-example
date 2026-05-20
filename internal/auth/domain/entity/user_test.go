package entity

import (
	"testing"
	"time"
)

const testUUID = "550e8400-e29b-41d4-a716-446655440000"

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
			if !tt.wantError && user.ID() == "" {
				t.Error("expected non-empty UUID for new user")
			}
		})
	}
}

func TestReconstructUser(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		id        string
		email     string
		password  string
		userName  string
		wantError bool
	}{
		{"valid", testUUID, "test@example.com", "hashed", "Test", false},
		{"empty id", "", "test@example.com", "hashed", "Test", true},
		{"empty email", testUUID, "", "hashed", "Test", true},
		{"empty password", testUUID, "test@example.com", "", "Test", true},
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
