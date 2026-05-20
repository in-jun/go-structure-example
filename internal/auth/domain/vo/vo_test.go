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
