package http

import (
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/auction/application/command"
	"github.com/in-jun/go-structure-example/internal/auction/application/query"
)

func TestToCreateResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	result := &command.CreateResult{
		ID:          "test-id",
		Title:       "Test",
		Description: "Desc",
		StartPrice:  1000,
		Status:      "draft",
		EndTime:     now.Add(24 * time.Hour),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	resp := toCreateResponse(result)

	if resp.ID != "test-id" {
		t.Errorf("ID = %q, want %q", resp.ID, "test-id")
	}
	if resp.Title != "Test" {
		t.Errorf("Title = %q, want %q", resp.Title, "Test")
	}
	if resp.StartPrice != 1000 {
		t.Errorf("StartPrice = %d, want 1000", resp.StartPrice)
	}
}

func TestToGetResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	result := &query.Result{
		ID:       "test-id",
		SellerID: "seller-id",
		Title:    "Auction",
		Status:   "open",
		EndTime:  now.Add(24 * time.Hour),
	}

	resp := toGetResponse(result)

	if resp.SellerID != "seller-id" {
		t.Errorf("SellerID = %q, want %q", resp.SellerID, "seller-id")
	}
}

func TestToListResponse(t *testing.T) {
	result := &query.ListResult{
		Auctions: []query.Result{
			{ID: "1", Title: "A1"},
			{ID: "2", Title: "A2"},
		},
		Total: 2,
	}

	resp := toListResponse(result)

	if resp.Total != 2 {
		t.Errorf("Total = %d, want 2", resp.Total)
	}
	if len(resp.Auctions) != 2 {
		t.Errorf("len(Auctions) = %d, want 2", len(resp.Auctions))
	}
}

func TestToEventHistoryResponse(t *testing.T) {
	now := time.Now()
	result := &query.EventHistoryResult{
		Events: []query.EventHistoryItem{
			{ID: 1, EventType: "auction.opened", Payload: []byte(`{}`), OccurredAt: now},
		},
	}

	resp := toEventHistoryResponse(result)

	if len(resp.Events) != 1 {
		t.Errorf("len(Events) = %d, want 1", len(resp.Events))
	}
	if resp.Events[0].EventType != "auction.opened" {
		t.Errorf("EventType = %q, want %q", resp.Events[0].EventType, "auction.opened")
	}
}
