package vo

import (
	"testing"
)

const validUUID = "550e8400-e29b-41d4-a716-446655440000"

func TestNewPaymentIDVO(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		wantError bool
	}{
		{"valid uuid", validUUID, false},
		{"invalid", "not-uuid", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := NewPaymentIDVO(tt.id)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", vo)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestNewAuctionIDVO(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		wantError bool
	}{
		{"valid uuid", validUUID, false},
		{"invalid", "bad", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := NewAuctionIDVO(tt.id)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", vo)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
