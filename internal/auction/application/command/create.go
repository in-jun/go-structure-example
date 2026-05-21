package command

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/auction/domain"
	"github.com/in-jun/go-structure-example/internal/auction/domain/entity"
	"github.com/in-jun/go-structure-example/internal/auction/domain/service"
	"github.com/in-jun/go-structure-example/internal/auction/domain/vo"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

type Create struct {
	UserID      string
	Title       string
	Description string
	StartPrice  int64
	EndTime     time.Time
}

type CreateResult struct {
	ID          string
	Title       string
	Description string
	StartPrice  int64
	Status      string
	EndTime     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateHandler struct {
	auctionRepo    domain.AuctionRepository
	eventPublisher domain.EventPublisher
	scheduler      *service.AuctionScheduler
	transactor     transaction.Transactor
}

func NewCreateHandler(auctionRepo domain.AuctionRepository, eventPublisher domain.EventPublisher, scheduler *service.AuctionScheduler, transactor transaction.Transactor) *CreateHandler {
	return &CreateHandler{auctionRepo: auctionRepo, eventPublisher: eventPublisher, scheduler: scheduler, transactor: transactor}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd Create) (*CreateResult, error) {
	sv, err := vo.NewSellerIDVO(cmd.UserID)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	cv, err := vo.NewCreateVO(cmd.Title, cmd.Description, cmd.StartPrice, cmd.EndTime)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	if err := h.scheduler.ValidateTiming(cv.EndTime); err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	auction, err := entity.NewAuction(sv.ID, cv.Title, cv.Description, cv.StartPrice, cv.EndTime)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	var result *CreateResult
	err = h.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := h.auctionRepo.Save(txCtx, auction); err != nil {
			return err
		}
		if err := h.eventPublisher.Publish(txCtx, auction.Events()...); err != nil {
			return err
		}
		auction.ClearEvents()

		result = &CreateResult{
			ID: auction.ID(), Title: auction.Title(), Description: auction.Description(),
			StartPrice: auction.StartPrice(), Status: auction.Status(), EndTime: auction.EndTime(),
			CreatedAt: auction.CreatedAt(), UpdatedAt: auction.UpdatedAt(),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
