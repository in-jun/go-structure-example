package nats

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/event"
	"github.com/in-jun/go-structure-example/internal/shared/observability"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type IdempotentHandler func(ctx context.Context, envelope *event.Envelope) error

func SubscribeIdempotent(
	nc *nats.Conn,
	subject string,
	serviceName string,
	dbGetter func(ctx context.Context) transaction.DBTX,
	transactor transaction.Transactor,
	handler IdempotentHandler,
) (*nats.Subscription, error) {
	return nc.Subscribe(subject, func(msg *nats.Msg) {
		var env event.Envelope
		if err := json.Unmarshal(msg.Data, &env); err != nil {
			slog.Error("failed to unmarshal event", "component", "nats", "subject", subject, "error", err)
			return
		}

		msgCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if len(env.TraceContext) > 0 {
			carrier := propagation.MapCarrier(env.TraceContext)
			msgCtx = otel.GetTextMapPropagator().Extract(msgCtx, carrier)
		}
		tracer := otel.Tracer("nats-consumer")
		msgCtx, span := tracer.Start(msgCtx, "consume:"+subject)
		defer span.End()

		err := transactor.WithinTransaction(msgCtx, func(txCtx context.Context) error {
			db := dbGetter(txCtx)

			var exists bool
			if err := db.QueryRowContext(txCtx,
				"SELECT EXISTS(SELECT 1 FROM processed_events WHERE event_id = $1)", env.ID,
			).Scan(&exists); err != nil {
				return err
			}
			if exists {
				slog.Info("skipping duplicate event", "component", "nats", "event_id", env.ID, "subject", subject)
				return nil
			}

			if err := handler(txCtx, &env); err != nil {
				return err
			}

			_, err := db.ExecContext(txCtx,
				"INSERT INTO processed_events (event_id) VALUES ($1)", env.ID)
			return err
		})

		if err != nil {
			span.RecordError(err)
			slog.Error("failed to handle event", "component", "nats", "subject", subject, "error", err)
			observability.NATSEventsConsumed.WithLabelValues(serviceName, subject, "error").Inc()
		} else {
			observability.NATSEventsConsumed.WithLabelValues(serviceName, subject, "success").Inc()
		}
	})
}
