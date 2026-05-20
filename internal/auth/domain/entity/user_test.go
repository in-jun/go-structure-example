package entity

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		userName string
		wantErr  bool
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
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if user != nil {
					t.Error("expected nil user")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user.Email() != tt.email {
					t.Errorf("expected email %q, got %q", tt.email, user.Email())
				}
				if user.Name() != tt.userName {
					t.Errorf("expected name %q, got %q", tt.userName, user.Name())
				}
				if user.ID() != 0 {
					t.Error("expected zero ID for new user before save")
				}
			}
		})
	}
}

func TestReconstructUser(t *testing.T) {
	now := time.Now()

	t.Run("valid", func(t *testing.T) {
		u, err := ReconstructUser(1, "test@example.com", "hashed", "Test", now, now)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if u.ID() != 1 {
			t.Errorf("expected ID 1, got %d", u.ID())
		}
		if u.Email() != "test@example.com" {
			t.Errorf("expected email test@example.com, got %q", u.Email())
		}
		if u.Name() != "Test" {
			t.Errorf("expected name Test, got %q", u.Name())
		}
	})

	t.Run("zero id", func(t *testing.T) {
		_, err := ReconstructUser(0, "test@example.com", "hashed", "Test", now, now)
		if err == nil {
			t.Error("expected error for zero id")
		}
	})

	t.Run("empty email", func(t *testing.T) {
		_, err := ReconstructUser(1, "", "hashed", "Test", now, now)
		if err == nil {
			t.Error("expected error for empty email")
		}
	})

	t.Run("empty password", func(t *testing.T) {
		_, err := ReconstructUser(1, "test@example.com", "", "Test", now, now)
		if err == nil {
			t.Error("expected error for empty password")
		}
	})
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
