package event

import (
	"testing"
	"time"
)

const testID = "test-auction-id"
const testSellerID = "test-seller-id"

func TestAuctionCreated_EventName(t *testing.T) {
	e := NewAuctionCreated(testID, testSellerID, "Title", 1000, time.Now().Add(time.Hour))
	if e.EventName() != "auction.created" {
		t.Errorf("EventName = %q, want auction.created", e.EventName())
	}
	if e.AggregateID() != testID {
		t.Errorf("AggregateID = %q, want %q", e.AggregateID(), testID)
	}
	if e.OccurredAt().IsZero() {
		t.Error("OccurredAt should not be zero")
	}
}

func TestAuctionOpened_EventName(t *testing.T) {
	e := NewAuctionOpened(testID, testSellerID, 1000, time.Now().Add(time.Hour))
	if e.EventName() != "auction.opened" {
		t.Errorf("EventName = %q, want auction.opened", e.EventName())
	}
	if e.AggregateID() != testID {
		t.Errorf("AggregateID = %q, want %q", e.AggregateID(), testID)
	}
}

func TestAuctionClosed_EventName(t *testing.T) {
	e := NewAuctionClosed(testID, testSellerID)
	if e.EventName() != "auction.closed" {
		t.Errorf("EventName = %q, want auction.closed", e.EventName())
	}
	if e.AggregateID() != testID {
		t.Errorf("AggregateID = %q, want %q", e.AggregateID(), testID)
	}
}

func TestAuctionSettled_EventName(t *testing.T) {
	e := NewAuctionSettled(testID)
	if e.EventName() != "auction.settled" {
		t.Errorf("EventName = %q, want auction.settled", e.EventName())
	}
	if e.AggregateID() != testID {
		t.Errorf("AggregateID = %q, want %q", e.AggregateID(), testID)
	}
}

func TestAuctionCancelled_EventName(t *testing.T) {
	e := NewAuctionCancelled(testID)
	if e.EventName() != "auction.cancelled" {
		t.Errorf("EventName = %q, want auction.cancelled", e.EventName())
	}
	if e.AggregateID() != testID {
		t.Errorf("AggregateID = %q, want %q", e.AggregateID(), testID)
	}
}
