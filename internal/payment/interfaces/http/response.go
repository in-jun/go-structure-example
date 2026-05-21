package http

import (
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
