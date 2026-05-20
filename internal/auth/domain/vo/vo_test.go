package vo

import "testing"

const testUUID = "550e8400-e29b-41d4-a716-446655440000"

func TestNewRegisterVO(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		password  string
		userName  string
		wantError bool
	}{
		{"valid", "test@example.com", "password123", "Test User", false},
		{"invalid email", "not-an-email", "password123", "Test", true},
		{"short password", "test@example.com", "12345", "Test", true},
		{"empty name", "test@example.com", "password123", "", true},
		{"empty email", "", "password123", "Test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := NewRegisterVO(tt.email, tt.password, tt.userName)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", vo)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestNewLoginVO(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		password  string
		wantError bool
	}{
		{"valid", "test@example.com", "password", false},
		{"empty email", "", "password", true},
		{"empty password", "test@example.com", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := NewLoginVO(tt.email, tt.password)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", vo)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestNewRefreshTokenVO(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		wantError bool
	}{
		{"valid", "some-token", false},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := NewRefreshTokenVO(tt.token)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", vo)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestNewTokenStringVO(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		wantError bool
	}{
		{"valid", "jwt-token", false},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := NewTokenStringVO(tt.token)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", vo)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestNewUserIDVO(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		wantError bool
	}{
		{"valid uuid", testUUID, false},
		{"empty", "", true},
		{"not a uuid", "not-a-uuid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := NewUserIDVO(tt.id)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", vo)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
