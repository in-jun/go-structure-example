package nats

import (
	"context"
	"log/slog"

	"github.com/in-jun/go-structure-example/internal/bid/application/command"
	sharedEvent "github.com/in-jun/go-structure-example/internal/shared/event"
	sharedNats "github.com/in-jun/go-structure-example/internal/shared/nats"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
	"github.com/nats-io/nats.go"
)

type Consumer struct {
	nc                     *nats.Conn
	determineWinnerHandler *command.DetermineWinnerHandler
	dbGetter               func(ctx context.Context) transaction.DBTX
	transactor             transaction.Transactor
	subs                   []*nats.Subscription
}

func NewConsumer(
	nc *nats.Conn,
	determineWinnerHandler *command.DetermineWinnerHandler,
	dbGetter func(ctx context.Context) transaction.DBTX,
	transactor transaction.Transactor,
) *Consumer {
	return &Consumer{
		nc: nc, determineWinnerHandler: determineWinnerHandler,
		dbGetter: dbGetter, transactor: transactor,
	}
}

func (c *Consumer) Start(_ context.Context) error {
	sub, err := sharedNats.SubscribeIdempotent(c.nc, "auction.closed", "bid", c.dbGetter, c.transactor,
		func(ctx context.Context, env *sharedEvent.Envelope) error {
			slog.Info("received auction.closed", "service", "bid", "auction_id", env.AggregateID)
			return c.determineWinnerHandler.Handle(ctx, command.DetermineWinner{AuctionID: env.AggregateID})
		})
	if err != nil {
		return err
	}
	c.subs = append(c.subs, sub)

	slog.Info("NATS consumer started", "service", "bid", "subjects", "auction.closed")
	return nil
}

func (c *Consumer) Stop() error {
	for _, sub := range c.subs {
		if err := sub.Drain(); err != nil {
			return err
		}
	}
	return nil
}
