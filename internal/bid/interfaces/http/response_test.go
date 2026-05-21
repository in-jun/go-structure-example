package http

import (
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/bid/application/command"
	"github.com/in-jun/go-structure-example/internal/bid/application/query"
)

func TestToPlaceBidResponse(t *testing.T) {
	result := &command.PlaceBidResult{
		ID:        "bid-1",
		AuctionID: "auction-1",
		BidderID:  "user-1",
		Amount:    2500,
	}

	resp := toPlaceBidResponse(result)

	if resp.ID != "bid-1" {
		t.Errorf("ID = %q, want %q", resp.ID, "bid-1")
	}
	if resp.AuctionID != "auction-1" {
		t.Errorf("AuctionID = %q, want %q", resp.AuctionID, "auction-1")
	}
	if resp.Amount != 2500 {
		t.Errorf("Amount = %d, want 2500", resp.Amount)
	}
}

func TestToGetResponse_Bid(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	result := &query.Result{
		ID:        "bid-1",
		AuctionID: "auction-1",
		BidderID:  "user-1",
		Amount:    3000,
		CreatedAt: now,
	}

	resp := toGetResponse(result)

	if resp.BidderID != "user-1" {
		t.Errorf("BidderID = %q, want %q", resp.BidderID, "user-1")
	}
	if resp.Amount != 3000 {
		t.Errorf("Amount = %d, want 3000", resp.Amount)
	}
	if !resp.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", resp.CreatedAt, now)
	}
}

func TestToListResponse_Bid(t *testing.T) {
	result := &query.ListResult{
		Bids: []query.Result{
			{ID: "b1", AuctionID: "a1", Amount: 1000},
			{ID: "b2", AuctionID: "a1", Amount: 2000},
			{ID: "b3", AuctionID: "a1", Amount: 1500},
		},
		Total: 3,
	}

	resp := toListResponse(result)

	if resp.Total != 3 {
		t.Errorf("Total = %d, want 3", resp.Total)
	}
	if len(resp.Bids) != 3 {
		t.Errorf("len(Bids) = %d, want 3", len(resp.Bids))
	}
	if resp.Bids[1].Amount != 2000 {
		t.Errorf("Bids[1].Amount = %d, want 2000", resp.Bids[1].Amount)
	}
}

func TestToEventHistoryResponse_Bid(t *testing.T) {
	now := time.Now()
	result := &query.EventHistoryResult{
		Events: []query.EventHistoryItem{
			{ID: 1, EventType: "bid.placed", Payload: []byte(`{"amount":1000}`), OccurredAt: now},
			{ID: 2, EventType: "bid.won", Payload: []byte(`{"amount":1000}`), OccurredAt: now},
		},
	}

	resp := toEventHistoryResponse(result)

	if len(resp.Events) != 2 {
		t.Errorf("len(Events) = %d, want 2", len(resp.Events))
	}
	if resp.Events[0].EventType != "bid.placed" {
		t.Errorf("Events[0].EventType = %q, want bid.placed", resp.Events[0].EventType)
	}
	if resp.Events[1].EventType != "bid.won" {
		t.Errorf("Events[1].EventType = %q, want bid.won", resp.Events[1].EventType)
	}
}

func TestToListResponse_Bid_Empty(t *testing.T) {
	result := &query.ListResult{Bids: []query.Result{}, Total: 0}
	resp := toListResponse(result)
	if resp.Total != 0 {
		t.Errorf("Total = %d, want 0", resp.Total)
	}
	if len(resp.Bids) != 0 {
		t.Errorf("len(Bids) = %d, want 0", len(resp.Bids))
	}
}
