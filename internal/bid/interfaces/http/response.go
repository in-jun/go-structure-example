package http

import (
	"encoding/json"
	"time"

	"github.com/in-jun/go-structure-example/internal/bid/application/command"
	"github.com/in-jun/go-structure-example/internal/bid/application/query"
)

type Response struct {
	ID        string    `json:"id"`
	AuctionID string    `json:"auction_id"`
	BidderID  string    `json:"bidder_id"`
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type ListResponse struct {
	Bids  []Response `json:"bids"`
	Total int64      `json:"total"`
}

func toPlaceBidResponse(r *command.PlaceBidResult) *Response {
	return &Response{
		ID:        r.ID,
		AuctionID: r.AuctionID,
		BidderID:  r.BidderID,
		Amount:    r.Amount,
		CreatedAt: r.CreatedAt,
	}
}

func toGetResponse(r *query.Result) *Response {
	return &Response{
		ID:        r.ID,
		AuctionID: r.AuctionID,
		BidderID:  r.BidderID,
		Amount:    r.Amount,
		CreatedAt: r.CreatedAt,
	}
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

func toListResponse(r *query.ListResult) *ListResponse {
	bids := make([]Response, len(r.Bids))
	for i, b := range r.Bids {
		bids[i] = Response{
			ID:        b.ID,
			AuctionID: b.AuctionID,
			BidderID:  b.BidderID,
			Amount:    b.Amount,
			CreatedAt: b.CreatedAt,
		}
	}
	return &ListResponse{Bids: bids, Total: r.Total}
}
