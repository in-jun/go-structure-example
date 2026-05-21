package grpc

import (
	"context"
	"log/slog"
	"runtime/debug"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/in-jun/go-structure-example/internal/shared/logging"
	"github.com/in-jun/go-structure-example/internal/shared/observability"
)

func recoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error("gRPC panic recovered", "method", info.FullMethod, "panic", r, "stack", string(debug.Stack()))
				err = status.Errorf(codes.Internal, "internal error")
			}
		}()
		return handler(ctx, req)
	}
}

func loggingInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start)

		st, _ := status.FromError(err)
		code := st.Code()

		logging.FromContext(ctx).Info("grpc",
			slog.String("method", info.FullMethod),
			slog.String("code", code.String()),
			slog.Duration("latency", duration),
		)

		observability.GRPCRequestsTotal.WithLabelValues("auction-service", info.FullMethod, code.String()).Inc()
		observability.GRPCRequestDuration.WithLabelValues("auction-service", info.FullMethod).Observe(duration.Seconds())

		return resp, err
	}
}
