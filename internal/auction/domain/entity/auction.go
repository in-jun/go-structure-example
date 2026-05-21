package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/in-jun/go-structure-example/internal/auction/domain/event"
)

const (
	StatusDraft     = "draft"
	StatusOpen      = "open"
	StatusClosed    = "closed"
	StatusSettled   = "settled"
	StatusCancelled = "cancelled"
)

var (
	ErrNotDraft       = errors.New("auction is not in draft status")
	ErrNotOpen        = errors.New("auction is not open")
	ErrNotClosed      = errors.New("auction is not closed")
	ErrNotOwner       = errors.New("not the auction owner")
	ErrCannotCancel   = errors.New("auction cannot be cancelled in current status")
	errInvalidInput   = errors.New("seller ID and title are required")
	errInvalidPrice   = errors.New("start price must be positive")
	errInvalidEndTime = errors.New("end time must be in the future")
)

type Auction struct {
	id          string
	sellerID    string
	title       string
	description string
	startPrice  int64
	status      string
	endTime     time.Time
	createdAt   time.Time
	updatedAt   time.Time

	events []event.Event
}

func NewAuction(sellerID, title, description string, startPrice int64, endTime time.Time) (*Auction, error) {
	if sellerID == "" || title == "" {
		return nil, errInvalidInput
	}
	if startPrice <= 0 {
		return nil, errInvalidPrice
	}
	if !endTime.After(time.Now()) {
		return nil, errInvalidEndTime
	}
	now := time.Now()
	a := &Auction{
		id:          uuid.New().String(),
		sellerID:    sellerID,
		title:       title,
		description: description,
		startPrice:  startPrice,
		status:      StatusDraft,
		endTime:     endTime,
		createdAt:   now,
		updatedAt:   now,
	}
	a.record(event.NewAuctionCreated(a.id, sellerID, title, startPrice, endTime))
	return a, nil
}

func ReconstructAuction(id, sellerID, title, description string, startPrice int64, status string, endTime, createdAt, updatedAt time.Time) *Auction {
	return &Auction{
		id: id, sellerID: sellerID, title: title, description: description,
		startPrice: startPrice, status: status, endTime: endTime,
		createdAt: createdAt, updatedAt: updatedAt,
	}
}

func (a *Auction) ID() string          { return a.id }
func (a *Auction) SellerID() string     { return a.sellerID }
func (a *Auction) Title() string        { return a.title }
func (a *Auction) Description() string  { return a.description }
func (a *Auction) StartPrice() int64    { return a.startPrice }
func (a *Auction) Status() string       { return a.status }
func (a *Auction) EndTime() time.Time   { return a.endTime }
func (a *Auction) CreatedAt() time.Time { return a.createdAt }
func (a *Auction) UpdatedAt() time.Time { return a.updatedAt }

func (a *Auction) IsOwnedBy(userID string) bool { return a.sellerID == userID }

func (a *Auction) Open() error {
	if a.status != StatusDraft {
		return ErrNotDraft
	}
	a.status = StatusOpen
	a.updatedAt = time.Now()
	a.record(event.NewAuctionOpened(a.id, a.sellerID, a.startPrice, a.endTime))
	return nil
}

func (a *Auction) Close() error {
	if a.status != StatusOpen {
		return ErrNotOpen
	}
	a.status = StatusClosed
	a.updatedAt = time.Now()
	a.record(event.NewAuctionClosed(a.id, a.sellerID))
	return nil
}

func (a *Auction) Settle() error {
	if a.status != StatusClosed {
		return ErrNotClosed
	}
	a.status = StatusSettled
	a.updatedAt = time.Now()
	a.record(event.NewAuctionSettled(a.id))
	return nil
}

func (a *Auction) Cancel() error {
	if a.status != StatusClosed && a.status != StatusDraft {
		return ErrCannotCancel
	}
	a.status = StatusCancelled
	a.updatedAt = time.Now()
	a.record(event.NewAuctionCancelled(a.id))
	return nil
}

func (a *Auction) Events() []event.Event { return a.events }
func (a *Auction) ClearEvents()          { a.events = nil }
func (a *Auction) record(e event.Event)  { a.events = append(a.events, e) }
