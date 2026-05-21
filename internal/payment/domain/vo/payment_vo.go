package vo

import "github.com/in-jun/go-structure-example/internal/shared/validation"

type PaymentIDVO struct {
	ID string
}

func NewPaymentIDVO(id string) (*PaymentIDVO, error) {
	parsed, err := validation.ParseUUID(id)
	if err != nil {
		return nil, err
	}
	return &PaymentIDVO{ID: parsed}, nil
}

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
