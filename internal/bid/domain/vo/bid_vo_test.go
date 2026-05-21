package vo

import (
	"testing"
)

const validUUID = "550e8400-e29b-41d4-a716-446655440000"

func TestNewAuctionIDVO(t *testing.T) {
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

func TestNewBidderIDVO(t *testing.T) {
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
			vo, err := NewBidderIDVO(tt.id)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", vo)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestNewPlaceBidVO(t *testing.T) {
	tests := []struct {
		name      string
		auctionID string
		bidderID  string
		amount    int64
		wantError bool
	}{
		{"valid", validUUID, validUUID, 1000, false},
		{"invalid auction id", "bad", validUUID, 1000, true},
		{"invalid bidder id", validUUID, "bad", 1000, true},
		{"zero amount", validUUID, validUUID, 0, true},
		{"negative amount", validUUID, validUUID, -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := NewPlaceBidVO(tt.auctionID, tt.bidderID, tt.amount)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", vo)
			}
			if !tt.wantError {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				if vo.Amount != tt.amount {
					t.Errorf("expected amount %d, got %d", tt.amount, vo.Amount)
				}
			}
		})
	}
}
