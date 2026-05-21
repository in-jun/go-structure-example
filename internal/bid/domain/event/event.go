package event

import (
	"time"

	sharedevent "github.com/in-jun/go-structure-example/internal/shared/event"
)

type Event = sharedevent.Event
type StoredEvent = sharedevent.StoredEvent

type BidPlaced struct {
	BidID     string    `json:"bid_id"`
	AuctionID string    `json:"auction_id"`
	BidderID  string    `json:"bidder_id"`
	Amount    int64     `json:"amount"`
	Timestamp time.Time `json:"occurred_at"`
}

func NewBidPlaced(bidID, auctionID, bidderID string, amount int64) BidPlaced {
	return BidPlaced{
		BidID: bidID, AuctionID: auctionID, BidderID: bidderID,
		Amount: amount, Timestamp: time.Now(),
	}
}

func (e BidPlaced) EventName() string    { return "bid.placed" }
func (e BidPlaced) AggregateID() string  { return e.AuctionID }
func (e BidPlaced) OccurredAt() time.Time { return e.Timestamp }

type BidWon struct {
	BidID     string    `json:"bid_id"`
	AuctionID string    `json:"auction_id"`
	WinnerID  string    `json:"winner_id"`
	Amount    int64     `json:"amount"`
	Timestamp time.Time `json:"occurred_at"`
}

func NewBidWon(bidID, auctionID, winnerID string, amount int64) BidWon {
	return BidWon{
		BidID: bidID, AuctionID: auctionID, WinnerID: winnerID,
		Amount: amount, Timestamp: time.Now(),
	}
}

func (e BidWon) EventName() string    { return "bid.won" }
func (e BidWon) AggregateID() string  { return e.AuctionID }
func (e BidWon) OccurredAt() time.Time { return e.Timestamp }
