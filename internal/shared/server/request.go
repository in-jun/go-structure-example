package server

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"
)

func Bind(r *http.Request, v any) error {
	ct := r.Header.Get("Content-Type")
	if ct != "" && !strings.HasPrefix(ct, "application/json") {
		return errUnsupportedContentType
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

var errUnsupportedContentType = errors.New("content-type must be application/json")

func PathParam(r *http.Request, name string) string {
	return r.PathValue(name)
}

func QueryDefault(r *http.Request, key, def string) string {
	if v := r.URL.Query().Get(key); v != "" {
		return v
	}
	return def
}

type userIDKey struct{}
type tokenJTIKey struct{}
type tokenIssuedAtKey struct{}

func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

func ContextWithTokenClaims(ctx context.Context, jti string, issuedAt int64) context.Context {
	ctx = context.WithValue(ctx, tokenJTIKey{}, jti)
	ctx = context.WithValue(ctx, tokenIssuedAtKey{}, issuedAt)
	return ctx
}

func UserID(r *http.Request) string {
	if id, ok := r.Context().Value(userIDKey{}).(string); ok {
		return id
	}
	return ""
}

func TokenJTI(r *http.Request) string {
	if jti, ok := r.Context().Value(tokenJTIKey{}).(string); ok {
		return jti
	}
	return ""
}

func TokenIssuedAt(r *http.Request) int64 {
	if iat, ok := r.Context().Value(tokenIssuedAtKey{}).(int64); ok {
		return iat
	}
	return 0
}

func ClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if i := strings.Index(xff, ","); i > 0 {
			return strings.TrimSpace(xff[:i])
		}
		return xff
	}
	if xri := r.Header.Get("X-Real-Ip"); xri != "" {
		return xri
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
