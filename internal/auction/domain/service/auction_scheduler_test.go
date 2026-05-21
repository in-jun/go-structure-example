package service

import (
	"errors"
	"testing"
	"time"
)

func TestAuctionScheduler_ValidateTiming(t *testing.T) {
	s := &AuctionScheduler{}

	tests := []struct {
		name    string
		endTime time.Time
		wantErr error
	}{
		{"valid 2 hours", time.Now().Add(2 * time.Hour), nil},
		{"valid 7 days", time.Now().Add(7 * 24 * time.Hour), nil},
		{"too short", time.Now().Add(30 * time.Minute), ErrDurationTooShort},
		{"too long", time.Now().Add(31 * 24 * time.Hour), ErrDurationTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.ValidateTiming(tt.endTime)
			if tt.wantErr == nil && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}
