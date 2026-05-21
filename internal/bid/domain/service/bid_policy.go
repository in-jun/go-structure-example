package service

import "errors"

var (
	ErrBidTooLow = errors.New("bid must be higher than current highest bid")
	ErrBelowMin  = errors.New("bid must be at least the start price")
)

const MinBidIncrement int64 = 100

type BidPolicy struct{}

func (p *BidPolicy) Validate(amount, startPrice int64, highestBid *int64) error {
	if highestBid == nil {
		if amount < startPrice {
			return ErrBelowMin
		}
		return nil
	}
	if amount < *highestBid+MinBidIncrement {
		return ErrBidTooLow
	}
	return nil
}
