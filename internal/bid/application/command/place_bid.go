package command

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/bid/domain"
	"github.com/in-jun/go-structure-example/internal/bid/domain/entity"
	"github.com/in-jun/go-structure-example/internal/bid/domain/service"
	"github.com/in-jun/go-structure-example/internal/bid/domain/vo"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

type PlaceBid struct {
	UserID    string
	AuctionID string
	Amount    int64
}

type PlaceBidResult struct {
	ID        string
	AuctionID string
	BidderID  string
	Amount    int64
}

type PlaceBidHandler struct {
	bidRepo        domain.BidRepository
	auctionClient  domain.AuctionClient
	bidPolicy      *service.BidPolicy
	eventPublisher domain.EventPublisher
	transactor     transaction.Transactor
}

func NewPlaceBidHandler(
	bidRepo domain.BidRepository,
	auctionClient domain.AuctionClient,
	bidPolicy *service.BidPolicy,
	eventPublisher domain.EventPublisher,
	transactor transaction.Transactor,
) *PlaceBidHandler {
	return &PlaceBidHandler{
		bidRepo: bidRepo, auctionClient: auctionClient,
		bidPolicy: bidPolicy, eventPublisher: eventPublisher,
		transactor: transactor,
	}
}

func (h *PlaceBidHandler) Handle(ctx context.Context, cmd PlaceBid) (*PlaceBidResult, error) {
	bv, err := vo.NewBidderIDVO(cmd.UserID)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}
	av, err := vo.NewAuctionIDVO(cmd.AuctionID)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}
	if cmd.Amount <= 0 {
		return nil, errors.BadRequest("Amount must be positive")
	}

	auction, err := h.auctionClient.GetAuction(ctx, av.ID)
	if err != nil {
		return nil, err
	}
	if auction.Status != domain.AuctionStatusOpen {
		return nil, errors.BadRequest("Auction is not open for bidding")
	}
	if auction.SellerID == bv.ID {
		return nil, errors.Forbidden("Cannot bid on your own auction")
	}

	var result *PlaceBidResult
	err = h.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		highest, err := h.bidRepo.FindHighestByAuctionID(txCtx, av.ID)
		if err != nil {
			return err
		}

		var highestAmount *int64
		if highest != nil {
			amt := highest.Amount()
			highestAmount = &amt
		}

		if err := h.bidPolicy.Validate(cmd.Amount, auction.StartPrice, highestAmount); err != nil {
			return errors.BadRequest(err.Error())
		}

		bid, err := entity.NewBid(av.ID, bv.ID, cmd.Amount)
		if err != nil {
			return errors.BadRequest(err.Error())
		}

		if err := h.bidRepo.Save(txCtx, bid); err != nil {
			return err
		}

		if err := h.eventPublisher.Publish(txCtx, bid.Events()...); err != nil {
			return err
		}
		bid.ClearEvents()

		result = &PlaceBidResult{
			ID: bid.ID(), AuctionID: bid.AuctionID(),
			BidderID: bid.BidderID(), Amount: bid.Amount(),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
