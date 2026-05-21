package service

import (
	"errors"
	"time"
)

var (
	ErrDurationTooShort = errors.New("auction duration must be at least 1 hour")
	ErrDurationTooLong  = errors.New("auction duration must not exceed 30 days")
)

const (
	MinAuctionDuration = 1 * time.Hour
	MaxAuctionDuration = 30 * 24 * time.Hour
)

type AuctionScheduler struct{}

func (s *AuctionScheduler) ValidateTiming(endTime time.Time) error {
	duration := time.Until(endTime)
	if duration < MinAuctionDuration {
		return ErrDurationTooShort
	}
	if duration > MaxAuctionDuration {
		return ErrDurationTooLong
	}
	return nil
}
