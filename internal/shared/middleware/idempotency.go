package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/in-jun/go-structure-example/internal/shared/server"
)

type cachedResponse struct {
	StatusCode  int             `json:"status_code"`
	ContentType string          `json:"content_type"`
	Body        json.RawMessage `json:"body"`
}

const idempotencyTTL = 24 * time.Hour

func Idempotency(redisClient *redis.Client) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("Idempotency-Key")
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			userID := server.UserID(r)
			redisKey := "idempotency:" + userID + ":" + key

			cached, err := redisClient.Get(r.Context(), redisKey).Bytes()
			if err == nil {
				var resp cachedResponse
				if json.Unmarshal(cached, &resp) == nil {
					w.Header().Set("Content-Type", resp.ContentType)
					w.WriteHeader(resp.StatusCode)
					if _, err := w.Write(resp.Body); err != nil {
						slog.Warn("failed to write cached response body", "error", err)
					}
					return
				}
			}

			lockKey := redisKey + ":lock"
			locked, _ := redisClient.SetNX(r.Context(), lockKey, "1", 30*time.Second).Result()
			if !locked {
				server.Error(w, http.StatusConflict, "A request with this idempotency key is already being processed")
				return
			}
			defer redisClient.Del(context.Background(), lockKey)

			rw := server.EnsureResponseWriter(w)
			next.ServeHTTP(rw, r)

			if rw.StatusCode >= 200 && rw.StatusCode < 300 {
				resp := cachedResponse{
					StatusCode:  rw.StatusCode,
					ContentType: rw.Header().Get("Content-Type"),
					Body:        rw.Body.Bytes(),
				}
				if data, err := json.Marshal(resp); err == nil {
					redisClient.Set(r.Context(), redisKey, data, idempotencyTTL)
				}
			}
		})
	}
}
