package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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

func RateLimit(client *redis.Client, burst int) gin.HandlerFunc {
	windowMs := int(time.Second / time.Millisecond)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now().Unix()
		key := fmt.Sprintf("ratelimit:%s:%d", ip, now)

		allowed, err := rateLimitScript.Run(c.Request.Context(), client, []string{key}, burst, windowMs).Int()
		if err != nil {
			c.Next()
			return
		}

		if allowed == 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"message": "Too many requests"})
			return
		}

		c.Next()
	}
}
