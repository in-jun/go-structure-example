package event

import (
	"context"
	"encoding/json"

	"github.com/in-jun/go-structure-example/internal/payment/domain"
	"github.com/in-jun/go-structure-example/internal/payment/domain/event"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

var _ domain.EventPublisher = (*pgPublisher)(nil)

type pgPublisher struct {
	dbGetter func(ctx context.Context) transaction.DBTX
}

func NewPublisher(dbGetter func(ctx context.Context) transaction.DBTX) *pgPublisher {
	return &pgPublisher{dbGetter: dbGetter}
}

func (p *pgPublisher) PublishWithIDs(ctx context.Context, events ...event.Event) ([]int64, error) {
	db := p.dbGetter(ctx)
	var ids []int64
	for _, e := range events {
		payload, err := json.Marshal(e)
		if err != nil {
			return nil, errors.Internal("Failed to serialize event")
		}
		var id int64
		err = db.QueryRowContext(ctx,
			"INSERT INTO domain_events (aggregate_type, aggregate_id, event_type, payload, occurred_at, published) VALUES ($1, $2, $3, $4, $5, FALSE) RETURNING id",
			"payment", e.AggregateID(), e.EventName(), payload, e.OccurredAt(),
		).Scan(&id)
		if err != nil {
			return nil, errors.Internal("Failed to store domain event")
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (p *pgPublisher) Publish(ctx context.Context, events ...event.Event) error {
	_, err := p.PublishWithIDs(ctx, events...)
	return err
}
