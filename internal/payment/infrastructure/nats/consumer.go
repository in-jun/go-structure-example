package nats

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/in-jun/go-structure-example/internal/payment/application/command"
	sharedEvent "github.com/in-jun/go-structure-example/internal/shared/event"
	sharedNats "github.com/in-jun/go-structure-example/internal/shared/nats"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
	"github.com/nats-io/nats.go"
)

type Consumer struct {
	nc                   *nats.Conn
	createPaymentHandler *command.CreatePaymentHandler
	dbGetter             func(ctx context.Context) transaction.DBTX
	transactor           transaction.Transactor
	subs                 []*nats.Subscription
}

func NewConsumer(
	nc *nats.Conn,
	createPaymentHandler *command.CreatePaymentHandler,
	dbGetter func(ctx context.Context) transaction.DBTX,
	transactor transaction.Transactor,
) *Consumer {
	return &Consumer{
		nc: nc, createPaymentHandler: createPaymentHandler,
		dbGetter: dbGetter, transactor: transactor,
	}
}

type bidWonEvent struct {
	AuctionID string `json:"auction_id"`
	WinnerID  string `json:"winner_id"`
	Amount    int64  `json:"amount"`
}

func (c *Consumer) Start(_ context.Context) error {
	sub, err := sharedNats.SubscribeIdempotent(c.nc, "bid.won", c.dbGetter, c.transactor,
		func(ctx context.Context, env *sharedEvent.Envelope) error {
			var be bidWonEvent
			if err := json.Unmarshal(env.Payload, &be); err != nil {
				return err
			}
			slog.Info("received bid.won", "service", "payment", "auction_id", be.AuctionID, "winner_id", be.WinnerID, "amount", be.Amount)
			_, err := c.createPaymentHandler.Handle(ctx, command.CreatePayment{
				AuctionID: be.AuctionID,
				WinnerID:  be.WinnerID,
				Amount:    be.Amount,
			})
			return err
		})
	if err != nil {
		return err
	}
	c.subs = append(c.subs, sub)

	slog.Info("NATS consumer started", "service", "payment", "subjects", "bid.won")
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
