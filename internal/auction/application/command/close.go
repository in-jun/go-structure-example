package command

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/auction/domain"
	"github.com/in-jun/go-structure-example/internal/auction/domain/vo"
	"github.com/in-jun/go-structure-example/internal/shared/query"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

type Close struct {
	UserID    string
	AuctionID string
}

type CloseHandler struct {
	auctionRepo    domain.AuctionRepository
	eventPublisher domain.EventPublisher
	transactor     transaction.Transactor
}

func NewCloseHandler(auctionRepo domain.AuctionRepository, eventPublisher domain.EventPublisher, transactor transaction.Transactor) *CloseHandler {
	return &CloseHandler{auctionRepo: auctionRepo, eventPublisher: eventPublisher, transactor: transactor}
}

func (h *CloseHandler) Handle(ctx context.Context, cmd Close) error {
	sv, err := vo.NewSellerIDVO(cmd.UserID)
	if err != nil {
		return errors.BadRequest(err.Error())
	}
	av, err := vo.NewAuctionIDVO(cmd.AuctionID)
	if err != nil {
		return errors.BadRequest(err.Error())
	}

	return h.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		auction, err := h.auctionRepo.FindByID(txCtx, av.ID, query.ForUpdate())
		if err != nil {
			return err
		}
		if auction == nil {
			return errors.NotFound("Auction not found")
		}
		if !auction.IsOwnedBy(sv.ID) {
			return errors.Forbidden("Not authorized")
		}

		if err := auction.Close(); err != nil {
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
