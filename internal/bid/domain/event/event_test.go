package event

import "testing"

const testBidID = "test-bid-id"
const testAuctionID = "test-auction-id"
const testBidderID = "test-bidder-id"

func TestBidPlaced_EventName(t *testing.T) {
	e := NewBidPlaced(testBidID, testAuctionID, testBidderID, 1000)
	if e.EventName() != "bid.placed" {
		t.Errorf("EventName = %q, want bid.placed", e.EventName())
	}
	// AggregateID for bid.placed is the auction ID (for outbox routing per auction)
	if e.AggregateID() != testAuctionID {
		t.Errorf("AggregateID = %q, want %q", e.AggregateID(), testAuctionID)
	}
	if e.OccurredAt().IsZero() {
		t.Error("OccurredAt should not be zero")
	}
}

func TestBidWon_EventName(t *testing.T) {
	e := NewBidWon(testBidID, testAuctionID, testBidderID, 1000)
	if e.EventName() != "bid.won" {
		t.Errorf("EventName = %q, want bid.won", e.EventName())
	}
	if e.AggregateID() != testAuctionID {
		t.Errorf("AggregateID = %q, want %q", e.AggregateID(), testAuctionID)
	}
}
