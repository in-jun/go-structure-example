package vo

import "testing"

func TestNewUpdateProfileVO(t *testing.T) {
	tests := []struct {
		name      string
		userName  string
		wantError bool
	}{
		{"valid", "Test User", false},
		{"empty name", "", true},
		{"too long", string(make([]byte, 101)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := NewUpdateProfileVO(tt.userName)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", vo)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestNewUpdatePasswordVO(t *testing.T) {
	tests := []struct {
		name      string
		current   string
		newPwd    string
		wantError bool
	}{
		{"valid", "current123", "newpass123", false},
		{"empty current", "", "newpass123", true},
		{"empty new", "current123", "", true},
		{"short new", "current123", "12345", true},
		{"same password", "current123", "current123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := NewUpdatePasswordVO(tt.current, tt.newPwd)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", vo)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
