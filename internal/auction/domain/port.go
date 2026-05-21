package domain

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/auction/domain/entity"
	"github.com/in-jun/go-structure-example/internal/auction/domain/event"
	"github.com/in-jun/go-structure-example/internal/shared/query"
)

type AuctionRepository interface {
	Save(ctx context.Context, auction *entity.Auction) error
	FindByID(ctx context.Context, id string, opts ...query.Option) (*entity.Auction, error)
	FindAll(ctx context.Context, page, limit int) ([]*entity.Auction, int64, error)
	Update(ctx context.Context, auction *entity.Auction) error
}

type EventPublisher interface {
	Publish(ctx context.Context, events ...event.Event) error
}

type EventReader interface {
	FindByAuctionID(ctx context.Context, auctionID string) ([]event.StoredEvent, error)
}
