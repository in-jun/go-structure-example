package http

import (
	"encoding/json"
	"time"

	"github.com/in-jun/go-structure-example/internal/auction/application/command"
	"github.com/in-jun/go-structure-example/internal/auction/application/query"
)

type Response struct {
	ID          string    `json:"id"`
	SellerID    string    `json:"seller_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartPrice  int64     `json:"start_price"`
	Status      string    `json:"status"`
	EndTime     time.Time `json:"end_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListResponse struct {
	Auctions []Response `json:"auctions"`
	Total    int64      `json:"total"`
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

func toCreateResponse(r *command.CreateResult) *Response {
	return &Response{
		ID:          r.ID,
		SellerID:    r.SellerID,
		Title:       r.Title,
		Description: r.Description,
		StartPrice:  r.StartPrice,
		Status:      r.Status,
		EndTime:     r.EndTime,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func toGetResponse(r *query.Result) *Response {
	return &Response{
		ID:          r.ID,
		SellerID:    r.SellerID,
		Title:       r.Title,
		Description: r.Description,
		StartPrice:  r.StartPrice,
		Status:      r.Status,
		EndTime:     r.EndTime,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func toListResponse(r *query.ListResult) *ListResponse {
	auctions := make([]Response, len(r.Auctions))
	for i, a := range r.Auctions {
		auctions[i] = Response{
			ID:          a.ID,
			SellerID:    a.SellerID,
			Title:       a.Title,
			Description: a.Description,
			StartPrice:  a.StartPrice,
			Status:      a.Status,
			EndTime:     a.EndTime,
			CreatedAt:   a.CreatedAt,
			UpdatedAt:   a.UpdatedAt,
		}
	}
	return &ListResponse{Auctions: auctions, Total: r.Total}
}

func toEventHistoryResponse(r *query.EventHistoryResult) *EventHistoryResponse {
	events := make([]EventResponse, len(r.Events))
	for i, e := range r.Events {
		events[i] = EventResponse{
			ID:         e.ID,
			EventType:  e.EventType,
			Payload:    e.Payload,
			OccurredAt: e.OccurredAt,
		}
	}
	return &EventHistoryResponse{Events: events}
}
