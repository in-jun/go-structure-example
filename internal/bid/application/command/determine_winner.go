package command

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/bid/domain"
	"github.com/in-jun/go-structure-example/internal/bid/domain/event"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

type DetermineWinner struct {
	AuctionID string
}

type DetermineWinnerHandler struct {
	bidRepo        domain.BidRepository
	eventPublisher domain.EventPublisher
	transactor     transaction.Transactor
}

func NewDetermineWinnerHandler(bidRepo domain.BidRepository, eventPublisher domain.EventPublisher, transactor transaction.Transactor) *DetermineWinnerHandler {
	return &DetermineWinnerHandler{bidRepo: bidRepo, eventPublisher: eventPublisher, transactor: transactor}
}

func (h *DetermineWinnerHandler) Handle(ctx context.Context, cmd DetermineWinner) error {
	return h.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		highest, err := h.bidRepo.FindHighestByAuctionID(txCtx, cmd.AuctionID)
		if err != nil {
			return err
		}
		if highest == nil {
			return errors.NotFound("No bids found for auction")
		}

		evt := event.NewBidWon(highest.ID(), highest.AuctionID(), highest.BidderID(), highest.Amount())
		return h.eventPublisher.Publish(txCtx, evt)
	})
}
