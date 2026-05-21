package query

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/auction/domain"
)

type List struct {
	Page  int
	Limit int
}

type ListResult struct {
	Auctions []Result
	Total    int64
}

type ListHandler struct {
	auctionRepo domain.AuctionRepository
}

func NewListHandler(auctionRepo domain.AuctionRepository) *ListHandler {
	return &ListHandler{auctionRepo: auctionRepo}
}

func (h *ListHandler) Handle(ctx context.Context, qry List) (*ListResult, error) {
	auctions, total, err := h.auctionRepo.FindAll(ctx, qry.Page, qry.Limit)
	if err != nil {
		return nil, err
	}

	results := make([]Result, len(auctions))
	for i, a := range auctions {
		results[i] = Result{
			ID: a.ID(), SellerID: a.SellerID(), Title: a.Title(),
			Description: a.Description(), StartPrice: a.StartPrice(),
			Status: a.Status(), EndTime: a.EndTime(),
			CreatedAt: a.CreatedAt(), UpdatedAt: a.UpdatedAt(),
		}
	}

	return &ListResult{Auctions: results, Total: total}, nil
}
