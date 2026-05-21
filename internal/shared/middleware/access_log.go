package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/logging"
	"github.com/in-jun/go-structure-example/internal/shared/server"
)

func AccessLog() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			path := r.URL.Path

			rw := server.EnsureResponseWriter(w)
			next.ServeHTTP(rw, r)

			logging.FromContext(r.Context()).Info("request",
				slog.String("method", r.Method),
				slog.String("path", path),
				slog.Int("status", rw.StatusCode),
				slog.Duration("latency", time.Since(start)),
				slog.String("ip", server.ClientIP(r)),
			)
		})
	}
}
