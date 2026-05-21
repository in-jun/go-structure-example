package pg

import (
	"context"
	"database/sql"
	stderrors "errors"
	"time"

	"github.com/in-jun/go-structure-example/internal/payment/domain"
	"github.com/in-jun/go-structure-example/internal/payment/domain/entity"
	"github.com/in-jun/go-structure-example/internal/shared/query"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

var _ domain.PaymentRepository = (*paymentRepository)(nil)

type paymentRepository struct {
	dbGetter func(ctx context.Context) transaction.DBTX
}

func NewPaymentRepository(dbGetter func(ctx context.Context) transaction.DBTX) domain.PaymentRepository {
	return &paymentRepository{dbGetter: dbGetter}
}

func (r *paymentRepository) Save(ctx context.Context, payment *entity.Payment) error {
	db := r.dbGetter(ctx)
	_, err := db.ExecContext(ctx,
		"INSERT INTO payments (id, auction_id, winner_id, amount, status) VALUES ($1, $2, $3, $4, $5)",
		payment.ID(), payment.AuctionID(), payment.WinnerID(), payment.Amount(), payment.Status(),
	)
	if err != nil {
		return errors.Internal("Failed to create payment")
	}
	return nil
}

func (r *paymentRepository) FindByID(ctx context.Context, id string, opts ...query.Option) (*entity.Payment, error) {
	cfg := query.ApplyOptions(opts)
	db := r.dbGetter(ctx)
	var pid, auctionID, winnerID, status string
	var amount int64
	var createdAt, updatedAt time.Time

	query := "SELECT id, auction_id, winner_id, amount, status, created_at, updated_at FROM payments WHERE id = $1"
	if cfg.ForUpdate {
		query += " FOR UPDATE"
	}

	err := db.QueryRowContext(ctx, query, id).Scan(&pid, &auctionID, &winnerID, &amount, &status, &createdAt, &updatedAt)
	if stderrors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Internal("Failed to get payment")
	}

	return entity.ReconstructPayment(pid, auctionID, winnerID, amount, status, createdAt, updatedAt), nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *entity.Payment) error {
	db := r.dbGetter(ctx)
	result, err := db.ExecContext(ctx,
		"UPDATE payments SET status = $1, updated_at = $2 WHERE id = $3",
		payment.Status(), payment.UpdatedAt(), payment.ID(),
	)
	if err != nil {
		return errors.Internal("Failed to update payment")
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("Failed to get affected rows")
	}
	if rows == 0 {
		return errors.NotFound("Payment not found")
	}
	return nil
}
