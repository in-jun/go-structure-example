package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/in-jun/go-structure-example/internal/bid/domain/event"
)

var (
	errInvalidInput  = errors.New("auction ID and bidder ID are required")
	errInvalidAmount = errors.New("bid amount must be positive")
	ErrSelfBid       = errors.New("seller cannot bid on own auction")
)

type Bid struct {
	id        string
	auctionID string
	bidderID  string
	amount    int64
	createdAt time.Time

	events []event.Event
}

func NewBid(auctionID, bidderID string, amount int64) (*Bid, error) {
	if auctionID == "" || bidderID == "" {
		return nil, errInvalidInput
	}
	if amount <= 0 {
		return nil, errInvalidAmount
	}
	now := time.Now()
	bid := &Bid{
		id:        uuid.New().String(),
		auctionID: auctionID,
		bidderID:  bidderID,
		amount:    amount,
		createdAt: now,
	}
	bid.record(event.NewBidPlaced(bid.id, auctionID, bidderID, amount))
	return bid, nil
}

func ReconstructBid(id, auctionID, bidderID string, amount int64, createdAt time.Time) *Bid {
	return &Bid{
		id: id, auctionID: auctionID, bidderID: bidderID,
		amount: amount, createdAt: createdAt,
	}
}

func (b *Bid) ID() string          { return b.id }
func (b *Bid) AuctionID() string   { return b.auctionID }
func (b *Bid) BidderID() string    { return b.bidderID }
func (b *Bid) Amount() int64       { return b.amount }
func (b *Bid) CreatedAt() time.Time { return b.createdAt }

func (b *Bid) Events() []event.Event { return b.events }
func (b *Bid) ClearEvents()          { b.events = nil }
func (b *Bid) record(e event.Event)  { b.events = append(b.events, e) }
