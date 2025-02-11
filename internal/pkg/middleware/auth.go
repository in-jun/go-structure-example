package middleware

import (
	"strings"

	"github.com/in-jun/go-structure-example/internal/pkg/utils/errors"
	"github.com/in-jun/go-structure-example/internal/pkg/utils/jwt"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(errors.Unauthorized("Missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Error(errors.Unauthorized("Invalid authorization header"))
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(parts[1])
		if err != nil {
			c.Error(errors.Unauthorized("Invalid token"))
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
