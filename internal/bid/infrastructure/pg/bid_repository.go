package pg

import (
	"context"
	"database/sql"
	stderrors "errors"
	"log/slog"
	"time"

	"github.com/in-jun/go-structure-example/internal/bid/domain"
	"github.com/in-jun/go-structure-example/internal/bid/domain/entity"
	"github.com/in-jun/go-structure-example/internal/shared/query"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

var _ domain.BidRepository = (*bidRepository)(nil)

type bidRepository struct {
	dbGetter func(ctx context.Context) transaction.DBTX
}

func NewBidRepository(dbGetter func(ctx context.Context) transaction.DBTX) domain.BidRepository {
	return &bidRepository{dbGetter: dbGetter}
}

func (r *bidRepository) Save(ctx context.Context, bid *entity.Bid) error {
	db := r.dbGetter(ctx)
	_, err := db.ExecContext(ctx,
		"INSERT INTO bids (id, auction_id, bidder_id, amount) VALUES ($1, $2, $3, $4)",
		bid.ID(), bid.AuctionID(), bid.BidderID(), bid.Amount(),
	)
	if err != nil {
		return errors.Internal("Failed to create bid")
	}
	return nil
}

func (r *bidRepository) FindHighestByAuctionID(ctx context.Context, auctionID string, opts ...query.Option) (*entity.Bid, error) {
	cfg := query.ApplyOptions(opts)
	db := r.dbGetter(ctx)
	var id, aucID, bidderID string
	var amount int64
	var createdAt time.Time

	q := "SELECT id, auction_id, bidder_id, amount, created_at FROM bids WHERE auction_id = $1 ORDER BY amount DESC LIMIT 1"
	if cfg.ForUpdate {
		q += " FOR UPDATE"
	}

	err := db.QueryRowContext(ctx, q, auctionID).Scan(&id, &aucID, &bidderID, &amount, &createdAt)
	if stderrors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get highest bid")
	}

	return entity.ReconstructBid(id, aucID, bidderID, amount, createdAt), nil
}

func (r *bidRepository) FindByAuctionID(ctx context.Context, auctionID string, page, limit int) ([]*entity.Bid, int64, error) {
	db := r.dbGetter(ctx)
	offset := (page - 1) * limit

	rows, err := db.QueryContext(ctx,
		"SELECT id, auction_id, bidder_id, amount, created_at, COUNT(*) OVER() FROM bids WHERE auction_id = $1 ORDER BY amount DESC LIMIT $2 OFFSET $3",
		auctionID, limit, offset,
	)
	if err != nil {
		return nil, 0, errors.Internal("Failed to list bids")
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var bids []*entity.Bid
	var total int64
	for rows.Next() {
		var id, aucID, bidderID string
		var amount int64
		var createdAt time.Time
		if err := rows.Scan(&id, &aucID, &bidderID, &amount, &createdAt, &total); err != nil {
			return nil, 0, errors.Internal("Failed to scan bid")
		}
		bids = append(bids, entity.ReconstructBid(id, aucID, bidderID, amount, createdAt))
	}
	if err := rows.Err(); err != nil {
		return nil, 0, errors.Internal("Error iterating bids")
	}

	return bids, total, nil
}
