package nats

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/in-jun/go-structure-example/internal/shared/event"
	"github.com/in-jun/go-structure-example/internal/shared/observability"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func serviceFromSubject(subject string) string {
	if i := strings.IndexByte(subject, '.'); i > 0 {
		return subject[:i]
	}
	return subject
}

func Publish(nc *nats.Conn, subject string, envelope *event.Envelope) error {
	data, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	if err := nc.Publish(subject, data); err != nil {
		return err
	}
	observability.NATSEventsPublished.WithLabelValues(serviceFromSubject(subject), subject).Inc()
	return nil
}

func PublishWithContext(ctx context.Context, nc *nats.Conn, subject string, envelope *event.Envelope) error {
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	envelope.TraceContext = carrier
	return Publish(nc, subject, envelope)
}
