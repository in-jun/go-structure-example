package middleware

import (
	"context"
	"strings"

	"github.com/in-jun/go-structure-example/internal/pkg/utils/errors"
	"github.com/in-jun/go-structure-example/internal/pkg/utils/jwt"

	"github.com/gin-gonic/gin"
)

type BlacklistChecker interface {
	IsAccessTokenBlacklisted(ctx context.Context, jti string) (bool, error)
}

var blacklistChecker BlacklistChecker

func SetBlacklistChecker(checker BlacklistChecker) {
	blacklistChecker = checker
}

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

		if blacklistChecker != nil {
			blacklisted, err := blacklistChecker.IsAccessTokenBlacklisted(c.Request.Context(), claims.ID)
			if err != nil {
				c.Error(errors.Internal("Failed to verify token"))
				c.Abort()
				return
			}
			if blacklisted {
				c.Error(errors.Unauthorized("Token has been revoked"))
				c.Abort()
				return
			}
		}

		c.Set("user_id", claims.UserID)
		c.Set("jti", claims.ID)
		c.Next()
	}
}
