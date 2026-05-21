package pg

import (
	"context"
	"database/sql"
	stderrors "errors"
	"time"

	"github.com/in-jun/go-structure-example/internal/auction/domain"
	"github.com/in-jun/go-structure-example/internal/auction/domain/entity"
	"github.com/in-jun/go-structure-example/internal/shared/query"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

var _ domain.AuctionRepository = (*auctionRepository)(nil)

type auctionRepository struct {
	dbGetter func(ctx context.Context) transaction.DBTX
}

func NewAuctionRepository(dbGetter func(ctx context.Context) transaction.DBTX) domain.AuctionRepository {
	return &auctionRepository{dbGetter: dbGetter}
}

func (r *auctionRepository) Save(ctx context.Context, auction *entity.Auction) error {
	db := r.dbGetter(ctx)
	_, err := db.ExecContext(ctx,
		"INSERT INTO auctions (id, seller_id, title, description, start_price, status, end_time) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		auction.ID(), auction.SellerID(), auction.Title(), auction.Description(), auction.StartPrice(), auction.Status(), auction.EndTime(),
	)
	if err != nil {
		return errors.Internal("Failed to create auction")
	}
	return nil
}

func (r *auctionRepository) FindByID(ctx context.Context, id string, opts ...query.Option) (*entity.Auction, error) {
	cfg := query.ApplyOptions(opts)
	db := r.dbGetter(ctx)
	var aid, sellerID, title, description, status string
	var startPrice int64
	var endTime, createdAt, updatedAt time.Time

	query := "SELECT id, seller_id, title, description, start_price, status, end_time, created_at, updated_at FROM auctions WHERE id = $1"
	if cfg.ForUpdate {
		query += " FOR UPDATE"
	}

	err := db.QueryRowContext(ctx, query, id).Scan(&aid, &sellerID, &title, &description, &startPrice, &status, &endTime, &createdAt, &updatedAt)
	if stderrors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get auction")
	}

	return entity.ReconstructAuction(aid, sellerID, title, description, startPrice, status, endTime, createdAt, updatedAt), nil
}

func (r *auctionRepository) FindAll(ctx context.Context, page, limit int) ([]*entity.Auction, int64, error) {
	db := r.dbGetter(ctx)
	offset := (page - 1) * limit

	rows, err := db.QueryContext(ctx,
		"SELECT id, seller_id, title, description, start_price, status, end_time, created_at, updated_at, COUNT(*) OVER() FROM auctions ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit, offset,
	)
	if err != nil {
		return nil, 0, errors.Internal("Failed to list auctions")
	}
	defer rows.Close()

	var auctions []*entity.Auction
	var total int64
	for rows.Next() {
		var aid, sellerID, title, description, status string
		var startPrice int64
		var endTime, createdAt, updatedAt time.Time
		if err := rows.Scan(&aid, &sellerID, &title, &description, &startPrice, &status, &endTime, &createdAt, &updatedAt, &total); err != nil {
			return nil, 0, errors.Internal("Failed to scan auction")
		}
		auctions = append(auctions, entity.ReconstructAuction(aid, sellerID, title, description, startPrice, status, endTime, createdAt, updatedAt))
	}
	if err := rows.Err(); err != nil {
		return nil, 0, errors.Internal("Error iterating auctions")
	}

	return auctions, total, nil
}

func (r *auctionRepository) Update(ctx context.Context, auction *entity.Auction) error {
	db := r.dbGetter(ctx)
	result, err := db.ExecContext(ctx,
		"UPDATE auctions SET title = $1, description = $2, start_price = $3, status = $4, end_time = $5, updated_at = $6 WHERE id = $7",
		auction.Title(), auction.Description(), auction.StartPrice(), auction.Status(), auction.EndTime(), auction.UpdatedAt(), auction.ID(),
	)
	if err != nil {
		return errors.Internal("Failed to update auction")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("Failed to get affected rows")
	}
	if rows == 0 {
		return errors.NotFound("Auction not found")
	}
	return nil
}
