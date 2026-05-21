package nats

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/in-jun/go-structure-example/internal/auction/application/command"
	sharedEvent "github.com/in-jun/go-structure-example/internal/shared/event"
	sharedNats "github.com/in-jun/go-structure-example/internal/shared/nats"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
	"github.com/nats-io/nats.go"
)

type paymentEvent struct {
	AuctionID string `json:"auction_id"`
}

type Consumer struct {
	nc            *nats.Conn
	settleHandler *command.SettleHandler
	cancelHandler *command.CancelHandler
	dbGetter      func(ctx context.Context) transaction.DBTX
	transactor    transaction.Transactor
	subs          []*nats.Subscription
}

func NewConsumer(
	nc *nats.Conn,
	settleHandler *command.SettleHandler,
	cancelHandler *command.CancelHandler,
	dbGetter func(ctx context.Context) transaction.DBTX,
	transactor transaction.Transactor,
) *Consumer {
	return &Consumer{
		nc: nc, settleHandler: settleHandler, cancelHandler: cancelHandler,
		dbGetter: dbGetter, transactor: transactor,
	}
}

func (c *Consumer) Start(_ context.Context) error {
	sub1, err := sharedNats.SubscribeIdempotent(c.nc, "payment.completed", "auction", c.dbGetter, c.transactor,
		func(ctx context.Context, env *sharedEvent.Envelope) error {
			var pe paymentEvent
			if err := json.Unmarshal(env.Payload, &pe); err != nil {
				return err
			}
			slog.Info("received payment.completed", "service", "auction", "auction_id", pe.AuctionID)
			return c.settleHandler.Handle(ctx, command.Settle{AuctionID: pe.AuctionID})
		})
	if err != nil {
		return err
	}
	c.subs = append(c.subs, sub1)

	sub2, err := sharedNats.SubscribeIdempotent(c.nc, "payment.failed", "auction", c.dbGetter, c.transactor,
		func(ctx context.Context, env *sharedEvent.Envelope) error {
			var pe paymentEvent
			if err := json.Unmarshal(env.Payload, &pe); err != nil {
				return err
			}
			slog.Info("received payment.failed", "service", "auction", "auction_id", pe.AuctionID)
			return c.cancelHandler.Handle(ctx, command.Cancel{AuctionID: pe.AuctionID})
		})
	if err != nil {
		return err
	}
	c.subs = append(c.subs, sub2)

	slog.Info("NATS consumers started", "service", "auction", "subjects", "payment.completed, payment.failed")
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
