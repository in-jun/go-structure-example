package service

import (
	"errors"
	"testing"
)

func int64Ptr(v int64) *int64 { return &v }

func TestBidPolicy_Validate(t *testing.T) {
	p := &BidPolicy{}

	tests := []struct {
		name       string
		amount     int64
		startPrice int64
		highest    *int64
		wantErr    error
	}{
		{"first bid at start price", 1000, 1000, nil, nil},
		{"first bid above start price", 1500, 1000, nil, nil},
		{"first bid below start price", 500, 1000, nil, ErrBelowMin},
		{"bid above highest + increment", 1200, 1000, int64Ptr(1000), nil},
		{"bid exactly at highest + increment", 1100, 1000, int64Ptr(1000), nil},
		{"bid below highest + increment", 1050, 1000, int64Ptr(1000), ErrBidTooLow},
		{"bid equal to highest", 1000, 1000, int64Ptr(1000), ErrBidTooLow},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := p.Validate(tt.amount, tt.startPrice, tt.highest)
			if tt.wantErr == nil && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}
