package validation

import (
	"testing"
)

func TestParseUUID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "550e8400-e29b-41d4-a716-446655440000", false},
		{"valid v4", "6ba7b810-9dad-11d1-80b4-00c04fd430c8", false},
		{"invalid", "not-a-uuid", true},
		{"empty", "", true},
		{"partial", "550e8400-e29b", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseUUID(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result != tt.input {
					t.Errorf("expected %q, got %q", tt.input, result)
				}
			}
		})
	}
}
