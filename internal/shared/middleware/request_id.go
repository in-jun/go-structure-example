package middleware

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/in-jun/go-structure-example/internal/shared/logging"
)

func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Request-ID")
			if id == "" {
				id = uuid.New().String()
			}
			w.Header().Set("X-Request-ID", id)
			ctx := logging.WithRequestID(r.Context(), id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
