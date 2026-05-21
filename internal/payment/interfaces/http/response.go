package http

import (
	"encoding/json"
	"time"

	"github.com/in-jun/go-structure-example/internal/payment/application/query"
)

type Response struct {
	ID        string    `json:"id"`
	AuctionID string    `json:"auction_id"`
	WinnerID  string    `json:"winner_id"`
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EventResponse struct {
	ID         int64           `json:"id"`
	EventType  string          `json:"event_type"`
	Payload    json.RawMessage `json:"payload"`
	OccurredAt time.Time       `json:"occurred_at"`
}

type EventHistoryResponse struct {
	Events []EventResponse `json:"events"`
}

func toEventHistoryResponse(r *query.EventHistoryResult) *EventHistoryResponse {
	events := make([]EventResponse, len(r.Events))
	for i, e := range r.Events {
		events[i] = EventResponse{ID: e.ID, EventType: e.EventType, Payload: e.Payload, OccurredAt: e.OccurredAt}
	}
	return &EventHistoryResponse{Events: events}
}

func toGetResponse(r *query.Result) *Response {
	return &Response{
		ID:        r.ID,
		AuctionID: r.AuctionID,
		WinnerID:  r.WinnerID,
		Amount:    r.Amount,
		Status:    r.Status,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}
