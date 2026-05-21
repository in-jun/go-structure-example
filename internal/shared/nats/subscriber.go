package nats

import (
	"encoding/json"
	"log/slog"

	"github.com/in-jun/go-structure-example/internal/shared/event"
	"github.com/nats-io/nats.go"
)

func Subscribe(nc *nats.Conn, subject string, handler func(envelope *event.Envelope) error) (*nats.Subscription, error) {
	return nc.Subscribe(subject, func(msg *nats.Msg) {
		var env event.Envelope
		if err := json.Unmarshal(msg.Data, &env); err != nil {
			slog.Error("failed to unmarshal event", "component", "nats", "subject", subject, "error", err)
			return
		}
		if err := handler(&env); err != nil {
			slog.Error("failed to handle event", "component", "nats", "subject", subject, "error", err)
		}
	})
}
