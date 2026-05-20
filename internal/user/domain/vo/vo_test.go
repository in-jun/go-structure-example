package vo

import "testing"

func TestNewUpdateProfileVO(t *testing.T) {
	tests := []struct {
		name    string
		uname   string
		wantErr bool
	}{
		{"valid", "Test User", false},
		{"empty name", "", true},
		{"too long", string(make([]byte, 101)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewUpdateProfileVO(tt.uname)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if v.Name != tt.uname {
					t.Errorf("expected name %q, got %q", tt.uname, v.Name)
				}
			}
		})
	}
}

func TestNewUpdatePasswordVO(t *testing.T) {
	tests := []struct {
		name     string
		current  string
		newPwd   string
		wantErr  bool
	}{
		{"valid", "current123", "newpass123", false},
		{"empty current", "", "newpass123", true},
		{"empty new", "current123", "", true},
		{"short new", "current123", "12345", true},
		{"same password", "current123", "current123", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUpdatePasswordVO(tt.current, tt.newPwd)
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
