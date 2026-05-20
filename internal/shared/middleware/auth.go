package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

type ValidateTokenResult struct {
	UserID   string
	JTI      string
	IssuedAt int64
}

type TokenValidator func(ctx context.Context, tokenString string) (*ValidateTokenResult, error)

func Auth(validateToken TokenValidator) gin.HandlerFunc {
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

		result, err := validateToken(c.Request.Context(), parts[1])
		if err != nil {
			c.Error(errors.Unauthorized("Invalid token"))
			c.Abort()
			return
		}

		c.Set("user_id", result.UserID)
		c.Set("jti", result.JTI)
		c.Set("issued_at", result.IssuedAt)
		c.Next()
	}
}
