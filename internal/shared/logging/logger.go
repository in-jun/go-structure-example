package logging

import (
	"context"
	"log/slog"
	"os"
)

type requestIDKey struct{}

func Init(serviceName string) {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(handler).With(slog.String("service", serviceName))
	slog.SetDefault(logger)
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return ""
}

func FromContext(ctx context.Context) *slog.Logger {
	logger := slog.Default()
	if id := RequestIDFromContext(ctx); id != "" {
		logger = logger.With(slog.String("request_id", id))
	}
	return logger
}
