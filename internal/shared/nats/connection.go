package nats

import (
	"log/slog"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/config"
	"github.com/nats-io/nats.go"
)

func NewConnection() (*nats.Conn, error) {
	return nats.Connect(config.AppConfig.NATSURL,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			slog.Error("nats disconnected", "error", err)
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			slog.Info("nats reconnected")
		}),
	)
}
