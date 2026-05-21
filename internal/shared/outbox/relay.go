package outbox

import (
	"context"
	"database/sql"
	"log/slog"
	"strconv"
	"time"

	sharedEvent "github.com/in-jun/go-structure-example/internal/shared/event"
	sharedNats "github.com/in-jun/go-structure-example/internal/shared/nats"
	"github.com/in-jun/go-structure-example/internal/shared/observability"
	"github.com/nats-io/nats.go"
)

type Relay struct {
	db          *sql.DB
	nc          *nats.Conn
	serviceName string
}

func NewRelay(db *sql.DB, nc *nats.Conn, serviceName string) *Relay {
	return &Relay{db: db, nc: nc, serviceName: serviceName}
}

func (r *Relay) Start(ctx context.Context) {
	slog.Info("outbox relay started", "component", "outbox-relay", "service", r.serviceName)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("outbox relay stopped", "component", "outbox-relay", "service", r.serviceName)
			return
		case <-ticker.C:
			if err := r.publishBatch(ctx); err != nil {
				slog.Error("outbox relay batch error", "component", "outbox-relay", "service", r.serviceName, "error", err)
			}
		}
	}
}

func (r *Relay) publishBatch(ctx context.Context) error {
	var pending float64
	if err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM domain_events WHERE published = FALSE AND aggregate_type = $1", r.serviceName,
	).Scan(&pending); err == nil {
		observability.OutboxPendingGauge.WithLabelValues(r.serviceName).Set(pending)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	rows, err := tx.QueryContext(ctx,
		"SELECT id, event_type, aggregate_id, payload, occurred_at FROM domain_events WHERE published = FALSE ORDER BY id LIMIT 100 FOR UPDATE SKIP LOCKED",
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id          int64
			eventType   string
			aggregateID string
			payload     []byte
			occurredAt  time.Time
		)
		if err := rows.Scan(&id, &eventType, &aggregateID, &payload, &occurredAt); err != nil {
			return err
		}

		env := &sharedEvent.Envelope{
			ID:          strconv.FormatInt(id, 10),
			Type:        eventType,
			AggregateID: aggregateID,
			Payload:     payload,
			OccurredAt:  occurredAt,
		}

		if err := sharedNats.Publish(r.nc, eventType, env); err != nil {
			slog.Warn("outbox relay NATS publish failed", "component", "outbox-relay", "service", r.serviceName, "event_id", id, "error", err)
			continue
		}

		if _, err := tx.ExecContext(ctx,
			"UPDATE domain_events SET published = TRUE WHERE id = $1", id); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	return tx.Commit()
}
