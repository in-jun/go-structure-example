package query

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/bid/domain"
	"github.com/in-jun/go-structure-example/internal/bid/domain/vo"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

type ListBids struct {
	AuctionID string
	Page      int
	Limit     int
}

type ListResult struct {
	Bids  []Result
	Total int64
}

type ListBidsHandler struct {
	bidRepo domain.BidRepository
}

func NewListBidsHandler(bidRepo domain.BidRepository) *ListBidsHandler {
	return &ListBidsHandler{bidRepo: bidRepo}
}

func (h *ListBidsHandler) Handle(ctx context.Context, qry ListBids) (*ListResult, error) {
	av, err := vo.NewAuctionIDVO(qry.AuctionID)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	bids, total, err := h.bidRepo.FindByAuctionID(ctx, av.ID, qry.Page, qry.Limit)
	if err != nil {
		return nil, err
	}

	results := make([]Result, len(bids))
	for i, b := range bids {
		results[i] = Result{
			ID: b.ID(), AuctionID: b.AuctionID(),
			BidderID: b.BidderID(), Amount: b.Amount(),
			CreatedAt: b.CreatedAt(),
		}
	}

	return &ListResult{Bids: results, Total: total}, nil
}
