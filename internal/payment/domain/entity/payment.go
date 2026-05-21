package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/in-jun/go-structure-example/internal/payment/domain/event"
)

const (
	StatusPending   = "pending"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
	StatusRefunded  = "refunded"
)

var (
	errInvalidInput    = errors.New("auction ID and winner ID are required")
	ErrNotPending      = errors.New("payment is not in pending status")
	ErrAlreadyComplete = errors.New("payment is already completed")
	ErrNotCompleted    = errors.New("payment is not in completed status")
)

type Payment struct {
	id        string
	auctionID string
	winnerID  string
	amount    int64
	status    string
	createdAt time.Time
	updatedAt time.Time

	events []event.Event
}

func NewPayment(auctionID, winnerID string, amount int64) (*Payment, error) {
	if auctionID == "" || winnerID == "" || amount <= 0 {
		return nil, errInvalidInput
	}
	now := time.Now()
	p := &Payment{
		id:        uuid.New().String(),
		auctionID: auctionID,
		winnerID:  winnerID,
		amount:    amount,
		status:    StatusPending,
		createdAt: now,
		updatedAt: now,
	}
	p.record(event.NewPaymentCreated(p.id, auctionID, winnerID, amount))
	return p, nil
}

func ReconstructPayment(id, auctionID, winnerID string, amount int64, status string, createdAt, updatedAt time.Time) *Payment {
	return &Payment{
		id: id, auctionID: auctionID, winnerID: winnerID,
		amount: amount, status: status,
		createdAt: createdAt, updatedAt: updatedAt,
	}
}

func (p *Payment) ID() string          { return p.id }
func (p *Payment) AuctionID() string   { return p.auctionID }
func (p *Payment) WinnerID() string    { return p.winnerID }
func (p *Payment) Amount() int64       { return p.amount }
func (p *Payment) Status() string      { return p.status }
func (p *Payment) CreatedAt() time.Time { return p.createdAt }
func (p *Payment) UpdatedAt() time.Time { return p.updatedAt }

func (p *Payment) Complete() error {
	if p.status != StatusPending {
		return ErrNotPending
	}
	p.status = StatusCompleted
	p.updatedAt = time.Now()
	p.record(event.NewPaymentCompleted(p.id, p.auctionID, p.winnerID, p.amount))
	return nil
}

func (p *Payment) Fail(reason string) error {
	if p.status != StatusPending {
		return ErrNotPending
	}
	p.status = StatusFailed
	p.updatedAt = time.Now()
	p.record(event.NewPaymentFailed(p.id, p.auctionID, p.winnerID, p.amount, reason))
	return nil
}

func (p *Payment) Refund(reason string) error {
	if p.status != StatusCompleted {
		return ErrNotCompleted
	}
	p.status = StatusRefunded
	p.updatedAt = time.Now()
	p.record(event.NewPaymentRefunded(p.id, p.auctionID, p.winnerID, p.amount, reason))
	return nil
}

func (p *Payment) IsOwnedBy(userID string) bool { return p.winnerID == userID }

func (p *Payment) Events() []event.Event { return p.events }
func (p *Payment) ClearEvents()          { p.events = nil }
func (p *Payment) record(e event.Event)  { p.events = append(p.events, e) }
