package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/observability"
	"github.com/in-jun/go-structure-example/internal/shared/server"
)

func Metrics(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			path := r.Pattern
			if path == "" {
				path = r.URL.Path
			}

			rw := server.EnsureResponseWriter(w)
			next.ServeHTTP(rw, r)

			status := strconv.Itoa(rw.StatusCode)
			duration := time.Since(start).Seconds()

			observability.HTTPRequestsTotal.WithLabelValues(serviceName, r.Method, path, status).Inc()
			observability.HTTPRequestDuration.WithLabelValues(serviceName, r.Method, path).Observe(duration)
		})
	}
}
