package http

import (
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
