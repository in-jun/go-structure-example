package middleware

import (
	"net/http"

	"github.com/in-jun/go-structure-example/internal/shared/server"
)

func BodyLimit(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxBytes {
				server.Error(w, http.StatusRequestEntityTooLarge, "Request body too large")
				return
			}
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}
