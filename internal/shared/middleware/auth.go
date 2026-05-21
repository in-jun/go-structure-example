package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/server"
)

type ValidateTokenResult struct {
	UserID   string
	JTI      string
	IssuedAt int64
}

type TokenValidator func(ctx context.Context, tokenString string) (*ValidateTokenResult, error)

func Auth(validateToken TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				HandleError(w, errors.Unauthorized("Missing authorization header"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				HandleError(w, errors.Unauthorized("Invalid authorization header"))
				return
			}

			result, err := validateToken(r.Context(), parts[1])
			if err != nil {
				HandleError(w, errors.Unauthorized("Invalid token"))
				return
			}

			ctx := server.ContextWithUserID(r.Context(), result.UserID)
			ctx = server.ContextWithTokenClaims(ctx, result.JTI, result.IssuedAt)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GatewayAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Header.Get("X-User-ID")
			if userID == "" {
				HandleError(w, errors.Unauthorized("Missing X-User-ID header"))
				return
			}
			if _, err := uuid.Parse(userID); err != nil {
				HandleError(w, errors.Unauthorized("Invalid X-User-ID header"))
				return
			}
			ctx := server.ContextWithUserID(r.Context(), userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
