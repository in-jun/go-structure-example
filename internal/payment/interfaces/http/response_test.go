package http

import (
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/payment/application/query"
)

func TestToGetResponse_Payment(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	result := &query.Result{
		ID:        "pay-1",
		AuctionID: "auction-1",
		WinnerID:  "user-1",
		Amount:    9999,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := toGetResponse(result)

	if resp.ID != "pay-1" {
		t.Errorf("ID = %q, want %q", resp.ID, "pay-1")
	}
	if resp.WinnerID != "user-1" {
		t.Errorf("WinnerID = %q, want %q", resp.WinnerID, "user-1")
	}
	if resp.Amount != 9999 {
		t.Errorf("Amount = %d, want 9999", resp.Amount)
	}
	if resp.Status != "pending" {
		t.Errorf("Status = %q, want pending", resp.Status)
	}
	if !resp.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", resp.CreatedAt, now)
	}
	if !resp.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v, want %v", resp.UpdatedAt, now)
	}
}

func TestToEventHistoryResponse_Payment(t *testing.T) {
	now := time.Now()
	result := &query.EventHistoryResult{
		Events: []query.EventHistoryItem{
			{ID: 1, EventType: "payment.created", Payload: []byte(`{"amount":9999}`), OccurredAt: now},
			{ID: 2, EventType: "payment.completed", Payload: []byte(`{}`), OccurredAt: now},
		},
	}

	resp := toEventHistoryResponse(result)

	if len(resp.Events) != 2 {
		t.Errorf("len(Events) = %d, want 2", len(resp.Events))
	}
	if resp.Events[0].EventType != "payment.created" {
		t.Errorf("Events[0].EventType = %q, want payment.created", resp.Events[0].EventType)
	}
	if resp.Events[1].EventType != "payment.completed" {
		t.Errorf("Events[1].EventType = %q, want payment.completed", resp.Events[1].EventType)
	}
}

func TestToEventHistoryResponse_Payment_Empty(t *testing.T) {
	result := &query.EventHistoryResult{Events: []query.EventHistoryItem{}}
	resp := toEventHistoryResponse(result)
	if len(resp.Events) != 0 {
		t.Errorf("len(Events) = %d, want 0", len(resp.Events))
	}
}
