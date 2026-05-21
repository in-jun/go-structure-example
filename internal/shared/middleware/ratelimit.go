package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/in-jun/go-structure-example/internal/shared/server"
)

var rateLimitScript = redis.NewScript(`
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local current = redis.call('INCR', key)
if current == 1 then
    redis.call('PEXPIRE', key, window)
end
if current > limit then
    return 0
end
return 1
`)

func RateLimit(redisClient *redis.Client, rps float64, burst int) func(http.Handler) http.Handler {
	windowMs := int(time.Second / time.Millisecond)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := server.ClientIP(r)
			now := time.Now().Unix()
			key := fmt.Sprintf("ratelimit:%s:%d", ip, now)

			allowed, err := rateLimitScript.Run(r.Context(), redisClient, []string{key}, burst, windowMs).Int()
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			if allowed == 0 {
				server.Error(w, http.StatusTooManyRequests, "Too Many Requests")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
