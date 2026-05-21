package nats

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/auction/domain"
	"github.com/in-jun/go-structure-example/internal/auction/domain/event"
	sharedEvent "github.com/in-jun/go-structure-example/internal/shared/event"
	sharedNats "github.com/in-jun/go-structure-example/internal/shared/nats"
	"github.com/nats-io/nats.go"
)

var _ domain.EventPublisher = (*publisher)(nil)

type publisher struct {
	nc *nats.Conn
}

func NewPublisher(nc *nats.Conn) domain.EventPublisher {
	return &publisher{nc: nc}
}

func (p *publisher) Publish(_ context.Context, events ...event.Event) error {
	for _, e := range events {
		env, err := sharedEvent.NewEnvelope(e.EventName(), e.AggregateID(), e, e.OccurredAt())
		if err != nil {
			return err
		}
		if err := sharedNats.Publish(p.nc, e.EventName(), env); err != nil {
			return err
		}
	}
	return nil
}
