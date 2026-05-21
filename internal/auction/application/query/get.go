package query

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/auction/domain"
	"github.com/in-jun/go-structure-example/internal/auction/domain/vo"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

type Get struct {
	AuctionID string
}

type Result struct {
	ID          string
	SellerID    string
	Title       string
	Description string
	StartPrice  int64
	Status      string
	EndTime     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GetHandler struct {
	auctionRepo domain.AuctionRepository
}

func NewGetHandler(auctionRepo domain.AuctionRepository) *GetHandler {
	return &GetHandler{auctionRepo: auctionRepo}
}

func (h *GetHandler) Handle(ctx context.Context, qry Get) (*Result, error) {
	av, err := vo.NewAuctionIDVO(qry.AuctionID)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	auction, err := h.auctionRepo.FindByID(ctx, av.ID)
	if err != nil {
		return nil, err
	}
	if auction == nil {
		return nil, errors.NotFound("Auction not found")
	}

	return &Result{
		ID: auction.ID(), SellerID: auction.SellerID(), Title: auction.Title(),
		Description: auction.Description(), StartPrice: auction.StartPrice(),
		Status: auction.Status(), EndTime: auction.EndTime(),
		CreatedAt: auction.CreatedAt(), UpdatedAt: auction.UpdatedAt(),
	}, nil
}
