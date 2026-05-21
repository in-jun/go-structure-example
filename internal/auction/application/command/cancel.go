package command

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/auction/domain"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/query"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

type Cancel struct {
	AuctionID string
}

type CancelHandler struct {
	auctionRepo    domain.AuctionRepository
	eventPublisher domain.EventPublisher
	transactor     transaction.Transactor
}

func NewCancelHandler(auctionRepo domain.AuctionRepository, eventPublisher domain.EventPublisher, transactor transaction.Transactor) *CancelHandler {
	return &CancelHandler{auctionRepo: auctionRepo, eventPublisher: eventPublisher, transactor: transactor}
}

func (h *CancelHandler) Handle(ctx context.Context, cmd Cancel) error {
	return h.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		auction, err := h.auctionRepo.FindByID(txCtx, cmd.AuctionID, query.ForUpdate())
		if err != nil {
			return err
		}
		if auction == nil {
			return errors.NotFound("Auction not found")
		}

		if err := auction.Cancel(); err != nil {
			return errors.Conflict(err.Error())
		}

		if err := h.auctionRepo.Update(txCtx, auction); err != nil {
			return err
		}

		if err := h.eventPublisher.Publish(txCtx, auction.Events()...); err != nil {
			return err
		}
		auction.ClearEvents()
		return nil
	}, transaction.WithIsolation(transaction.Pessimistic))
}
