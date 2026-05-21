package event

import (
	"time"

	sharedevent "github.com/in-jun/go-structure-example/internal/shared/event"
)

type Event = sharedevent.Event
type StoredEvent = sharedevent.StoredEvent

type AuctionCreated struct {
	AuctionID   string    `json:"auction_id"`
	SellerID    string    `json:"seller_id"`
	Title       string    `json:"title"`
	StartPrice  int64     `json:"start_price"`
	EndTime     time.Time `json:"end_time"`
	Timestamp   time.Time `json:"occurred_at"`
}

func NewAuctionCreated(auctionID, sellerID, title string, startPrice int64, endTime time.Time) AuctionCreated {
	return AuctionCreated{
		AuctionID: auctionID, SellerID: sellerID, Title: title,
		StartPrice: startPrice, EndTime: endTime,
		Timestamp: time.Now(),
	}
}

func (e AuctionCreated) EventName() string     { return "auction.created" }
func (e AuctionCreated) AggregateID() string   { return e.AuctionID }
func (e AuctionCreated) OccurredAt() time.Time { return e.Timestamp }

type AuctionOpened struct {
	AuctionID  string    `json:"auction_id"`
	SellerID   string    `json:"seller_id"`
	StartPrice int64     `json:"start_price"`
	EndTime    time.Time `json:"end_time"`
	Timestamp  time.Time `json:"occurred_at"`
}

func NewAuctionOpened(auctionID, sellerID string, startPrice int64, endTime time.Time) AuctionOpened {
	return AuctionOpened{
		AuctionID: auctionID, SellerID: sellerID,
		StartPrice: startPrice, EndTime: endTime,
		Timestamp: time.Now(),
	}
}

func (e AuctionOpened) EventName() string     { return "auction.opened" }
func (e AuctionOpened) AggregateID() string   { return e.AuctionID }
func (e AuctionOpened) OccurredAt() time.Time { return e.Timestamp }

type AuctionClosed struct {
	AuctionID string    `json:"auction_id"`
	SellerID  string    `json:"seller_id"`
	Timestamp time.Time `json:"occurred_at"`
}

func NewAuctionClosed(auctionID, sellerID string) AuctionClosed {
	return AuctionClosed{AuctionID: auctionID, SellerID: sellerID, Timestamp: time.Now()}
}

func (e AuctionClosed) EventName() string     { return "auction.closed" }
func (e AuctionClosed) AggregateID() string   { return e.AuctionID }
func (e AuctionClosed) OccurredAt() time.Time { return e.Timestamp }

type AuctionSettled struct {
	AuctionID string    `json:"auction_id"`
	Timestamp time.Time `json:"occurred_at"`
}

func NewAuctionSettled(auctionID string) AuctionSettled {
	return AuctionSettled{AuctionID: auctionID, Timestamp: time.Now()}
}

func (e AuctionSettled) EventName() string     { return "auction.settled" }
func (e AuctionSettled) AggregateID() string   { return e.AuctionID }
func (e AuctionSettled) OccurredAt() time.Time { return e.Timestamp }

type AuctionCancelled struct {
	AuctionID string    `json:"auction_id"`
	Timestamp time.Time `json:"occurred_at"`
}

func NewAuctionCancelled(auctionID string) AuctionCancelled {
	return AuctionCancelled{AuctionID: auctionID, Timestamp: time.Now()}
}

func (e AuctionCancelled) EventName() string     { return "auction.cancelled" }
func (e AuctionCancelled) AggregateID() string   { return e.AuctionID }
func (e AuctionCancelled) OccurredAt() time.Time { return e.Timestamp }
