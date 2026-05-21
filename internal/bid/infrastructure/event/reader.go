package event

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/bid/domain"
	domainEvent "github.com/in-jun/go-structure-example/internal/bid/domain/event"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

var _ domain.EventReader = (*pgReader)(nil)

type pgReader struct {
	dbGetter func(ctx context.Context) transaction.DBTX
}

func NewReader(dbGetter func(ctx context.Context) transaction.DBTX) domain.EventReader {
	return &pgReader{dbGetter: dbGetter}
}

func (r *pgReader) FindByAuctionID(ctx context.Context, auctionID string) ([]domainEvent.StoredEvent, error) {
	db := r.dbGetter(ctx)
	rows, err := db.QueryContext(ctx,
		"SELECT id, event_type, payload, occurred_at FROM domain_events WHERE aggregate_type = 'bid' AND aggregate_id = $1 ORDER BY id",
		auctionID,
	)
	if err != nil {
		return nil, errors.Internal("Failed to query domain events")
	}
	defer rows.Close()

	var events []domainEvent.StoredEvent
	for rows.Next() {
		var e domainEvent.StoredEvent
		if err := rows.Scan(&e.ID, &e.EventType, &e.Payload, &e.OccurredAt); err != nil {
			return nil, errors.Internal("Failed to scan domain event")
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Internal("Error iterating domain events")
	}
	return events, nil
}
