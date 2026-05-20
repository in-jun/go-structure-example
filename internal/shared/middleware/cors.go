package middleware

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/in-jun/go-structure-example/internal/shared/config"
)

func CORS() gin.HandlerFunc {
	allowOrigins := config.AppConfig.CORSAllowOrigins
	if allowOrigins == "" {
		allowOrigins = "*"
	}

	origins := strings.Split(allowOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: allowOrigins != "*",
		MaxAge:           12 * time.Hour,
	})
}
