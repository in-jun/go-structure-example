package vo

import "testing"

func TestNewRegisterVO(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		uname    string
		wantErr  bool
	}{
		{"valid", "test@example.com", "password123", "Test", false},
		{"empty email", "", "password123", "Test", true},
		{"invalid email", "not-an-email", "password123", "Test", true},
		{"short password", "test@example.com", "12345", "Test", true},
		{"empty name", "test@example.com", "password123", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewRegisterVO(tt.email, tt.password, tt.uname)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if v.Email != tt.email {
					t.Errorf("expected email %q, got %q", tt.email, v.Email)
				}
			}
		})
	}
}

func TestNewLoginVO(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		wantErr  bool
	}{
		{"valid", "test@example.com", "password123", false},
		{"empty email", "", "password123", true},
		{"empty password", "test@example.com", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewLoginVO(tt.email, tt.password)
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestNewRefreshTokenVO(t *testing.T) {
	if _, err := NewRefreshTokenVO(""); err == nil {
		t.Error("expected error for empty token, got nil")
	}
	v, err := NewRefreshTokenVO("some-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Token != "some-token" {
		t.Errorf("expected 'some-token', got %q", v.Token)
	}
}

func TestNewTokenStringVO(t *testing.T) {
	if _, err := NewTokenStringVO(""); err == nil {
		t.Error("expected error for empty token, got nil")
	}
	v, err := NewTokenStringVO("bearer-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Token != "bearer-token" {
		t.Errorf("expected 'bearer-token', got %q", v.Token)
	}
}

func TestNewUserIDVO(t *testing.T) {
	if _, err := NewUserIDVO(0); err == nil {
		t.Error("expected error for zero ID, got nil")
	}
	v, err := NewUserIDVO(42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.ID != 42 {
		t.Errorf("expected ID 42, got %d", v.ID)
	}
}
