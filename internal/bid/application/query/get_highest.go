package query

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/bid/domain"
	"github.com/in-jun/go-structure-example/internal/bid/domain/vo"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

type GetHighest struct {
	AuctionID string
}

type Result struct {
	ID        string
	AuctionID string
	BidderID  string
	Amount    int64
	CreatedAt time.Time
}

type GetHighestHandler struct {
	bidRepo domain.BidRepository
}

func NewGetHighestHandler(bidRepo domain.BidRepository) *GetHighestHandler {
	return &GetHighestHandler{bidRepo: bidRepo}
}

func (h *GetHighestHandler) Handle(ctx context.Context, qry GetHighest) (*Result, error) {
	av, err := vo.NewAuctionIDVO(qry.AuctionID)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	bid, err := h.bidRepo.FindHighestByAuctionID(ctx, av.ID)
	if err != nil {
		return nil, err
	}
	if bid == nil {
		return nil, errors.NotFound("No bids found")
	}

	return &Result{
		ID: bid.ID(), AuctionID: bid.AuctionID(),
		BidderID: bid.BidderID(), Amount: bid.Amount(),
		CreatedAt: bid.CreatedAt(),
	}, nil
}
