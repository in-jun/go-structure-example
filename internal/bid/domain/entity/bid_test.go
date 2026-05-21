package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	testAuctionID = uuid.New().String()
	testBidderID  = uuid.New().String()
)

func TestNewBid(t *testing.T) {
	bid, err := NewBid(testAuctionID, testBidderID, 1000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := uuid.Parse(bid.ID()); err != nil {
		t.Errorf("expected valid UUID, got '%s'", bid.ID())
	}
	if bid.AuctionID() != testAuctionID {
		t.Errorf("expected AuctionID '%s', got '%s'", testAuctionID, bid.AuctionID())
	}
	if bid.BidderID() != testBidderID {
		t.Errorf("expected BidderID '%s', got '%s'", testBidderID, bid.BidderID())
	}
	if bid.Amount() != 1000 {
		t.Errorf("expected Amount 1000, got %d", bid.Amount())
	}
	if len(bid.Events()) != 1 {
		t.Errorf("expected 1 event (BidPlaced), got %d", len(bid.Events()))
	}
}

func TestNewBid_Invariants(t *testing.T) {
	tests := []struct {
		name      string
		auctionID string
		bidderID  string
		amount    int64
	}{
		{"empty auctionID", "", testBidderID, 1000},
		{"empty bidderID", testAuctionID, "", 1000},
		{"both empty", "", "", 1000},
		{"zero amount", testAuctionID, testBidderID, 0},
		{"negative amount", testAuctionID, testBidderID, -500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewBid(tt.auctionID, tt.bidderID, tt.amount)
			if err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}

func TestBid_ClearEvents(t *testing.T) {
	bid, _ := NewBid(testAuctionID, testBidderID, 1000)
	bid.ClearEvents()
	if len(bid.Events()) != 0 {
		t.Error("expected empty events after ClearEvents")
	}
}

func TestReconstructBid(t *testing.T) {
	id := uuid.New().String()
	now := time.Now()
	bid := ReconstructBid(id, testAuctionID, testBidderID, 500, now)

	if bid.ID() != id {
		t.Errorf("expected ID '%s', got '%s'", id, bid.ID())
	}
	if len(bid.Events()) != 0 {
		t.Error("reconstructed bid should have no events")
	}
}
