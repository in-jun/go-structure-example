package domain

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/bid/domain/entity"
	"github.com/in-jun/go-structure-example/internal/bid/domain/event"
	"github.com/in-jun/go-structure-example/internal/shared/query"
)

const AuctionStatusOpen = "open"

type BidRepository interface {
	Save(ctx context.Context, bid *entity.Bid) error
	FindHighestByAuctionID(ctx context.Context, auctionID string, opts ...query.Option) (*entity.Bid, error)
	FindByAuctionID(ctx context.Context, auctionID string, page, limit int) ([]*entity.Bid, int64, error)
}

type AuctionClient interface {
	GetAuction(ctx context.Context, auctionID string) (*AuctionInfo, error)
}

type AuctionInfo struct {
	ID         string
	SellerID   string
	StartPrice int64
	Status     string
}

type EventPublisher interface {
	Publish(ctx context.Context, events ...event.Event) error
}

type EventReader interface {
	FindByAuctionID(ctx context.Context, auctionID string) ([]event.StoredEvent, error)
}
