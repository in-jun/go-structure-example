package vo

import (
	"errors"

	"github.com/in-jun/go-structure-example/internal/shared/validation"
)

var errInvalidAmount = errors.New("amount must be positive")

type AuctionIDVO struct {
	ID string
}

func NewAuctionIDVO(id string) (*AuctionIDVO, error) {
	parsed, err := validation.ParseUUID(id)
	if err != nil {
		return nil, err
	}
	return &AuctionIDVO{ID: parsed}, nil
}

type BidderIDVO struct {
	ID string
}

func NewBidderIDVO(id string) (*BidderIDVO, error) {
	parsed, err := validation.ParseUUID(id)
	if err != nil {
		return nil, err
	}
	return &BidderIDVO{ID: parsed}, nil
}

type PlaceBidVO struct {
	AuctionID string
	BidderID  string
	Amount    int64
}

func NewPlaceBidVO(auctionID, bidderID string, amount int64) (*PlaceBidVO, error) {
	av, err := NewAuctionIDVO(auctionID)
	if err != nil {
		return nil, err
	}
	bv, err := NewBidderIDVO(bidderID)
	if err != nil {
		return nil, err
	}
	if amount <= 0 {
		return nil, errInvalidAmount
	}
	return &PlaceBidVO{AuctionID: av.ID, BidderID: bv.ID, Amount: amount}, nil
}
