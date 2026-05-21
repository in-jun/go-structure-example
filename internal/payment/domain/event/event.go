package event

import (
	"time"

	sharedevent "github.com/in-jun/go-structure-example/internal/shared/event"
)

type Event = sharedevent.Event
type StoredEvent = sharedevent.StoredEvent

type PaymentCreated struct {
	PaymentID string    `json:"payment_id"`
	AuctionID string    `json:"auction_id"`
	WinnerID  string    `json:"winner_id"`
	Amount    int64     `json:"amount"`
	Timestamp time.Time `json:"occurred_at"`
}

func NewPaymentCreated(paymentID, auctionID, winnerID string, amount int64) PaymentCreated {
	return PaymentCreated{
		PaymentID: paymentID, AuctionID: auctionID, WinnerID: winnerID,
		Amount: amount, Timestamp: time.Now(),
	}
}

func (e PaymentCreated) EventName() string     { return "payment.created" }
func (e PaymentCreated) AggregateID() string   { return e.PaymentID }
func (e PaymentCreated) OccurredAt() time.Time { return e.Timestamp }

type PaymentCompleted struct {
	PaymentID string    `json:"payment_id"`
	AuctionID string    `json:"auction_id"`
	WinnerID  string    `json:"winner_id"`
	Amount    int64     `json:"amount"`
	Timestamp time.Time `json:"occurred_at"`
}

func NewPaymentCompleted(paymentID, auctionID, winnerID string, amount int64) PaymentCompleted {
	return PaymentCompleted{
		PaymentID: paymentID, AuctionID: auctionID, WinnerID: winnerID,
		Amount: amount, Timestamp: time.Now(),
	}
}

func (e PaymentCompleted) EventName() string    { return "payment.completed" }
func (e PaymentCompleted) AggregateID() string  { return e.PaymentID }
func (e PaymentCompleted) OccurredAt() time.Time { return e.Timestamp }

type PaymentFailed struct {
	PaymentID string    `json:"payment_id"`
	AuctionID string    `json:"auction_id"`
	WinnerID  string    `json:"winner_id"`
	Amount    int64     `json:"amount"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"occurred_at"`
}

func NewPaymentFailed(paymentID, auctionID, winnerID string, amount int64, reason string) PaymentFailed {
	return PaymentFailed{
		PaymentID: paymentID, AuctionID: auctionID, WinnerID: winnerID,
		Amount: amount, Reason: reason, Timestamp: time.Now(),
	}
}

func (e PaymentFailed) EventName() string    { return "payment.failed" }
func (e PaymentFailed) AggregateID() string  { return e.PaymentID }
func (e PaymentFailed) OccurredAt() time.Time { return e.Timestamp }

type PaymentRefunded struct {
	PaymentID string    `json:"payment_id"`
	AuctionID string    `json:"auction_id"`
	WinnerID  string    `json:"winner_id"`
	Amount    int64     `json:"amount"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"occurred_at"`
}

func NewPaymentRefunded(paymentID, auctionID, winnerID string, amount int64, reason string) PaymentRefunded {
	return PaymentRefunded{
		PaymentID: paymentID, AuctionID: auctionID, WinnerID: winnerID,
		Amount: amount, Reason: reason, Timestamp: time.Now(),
	}
}

func (e PaymentRefunded) EventName() string     { return "payment.refunded" }
func (e PaymentRefunded) AggregateID() string   { return e.PaymentID }
func (e PaymentRefunded) OccurredAt() time.Time { return e.Timestamp }
