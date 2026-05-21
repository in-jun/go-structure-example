package event

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/in-jun/go-structure-example/internal/bid/domain"
	domainEvent "github.com/in-jun/go-structure-example/internal/bid/domain/event"
	sharedEvent "github.com/in-jun/go-structure-example/internal/shared/event"
	sharedNats "github.com/in-jun/go-structure-example/internal/shared/nats"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
	"github.com/nats-io/nats.go"
)

type compositePublisher struct {
	pgPub *pgPublisher
	nc    *nats.Conn
}

var _ domain.EventPublisher = (*compositePublisher)(nil)

func NewCompositePublisher(pgPub *pgPublisher, nc *nats.Conn) domain.EventPublisher {
	return &compositePublisher{pgPub: pgPub, nc: nc}
}

func (c *compositePublisher) Publish(ctx context.Context, events ...domainEvent.Event) error {
	ids, err := c.pgPub.PublishWithIDs(ctx, events...)
	if err != nil {
		return err
	}

	evts := make([]domainEvent.Event, len(events))
	copy(evts, events)
	eventIDs := make([]int64, len(ids))
	copy(eventIDs, ids)

	transaction.RegisterPostCommit(ctx, func() {
		for i, e := range evts {
			envID := strconv.FormatInt(eventIDs[i], 10)
			env, err := sharedEvent.NewEnvelopeWithID(envID, e.EventName(), e.AggregateID(), e, e.OccurredAt())
			if err != nil {
				slog.Error("failed to create envelope", "component", "outbox", "error", err)
				continue
			}
			if err := sharedNats.Publish(c.nc, e.EventName(), env); err != nil {
				slog.Warn("NATS publish failed, relay will retry", "component", "outbox", "event_id", eventIDs[i], "error", err)
			}
		}
	})
	return nil
}
