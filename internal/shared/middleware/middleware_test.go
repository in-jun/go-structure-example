package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestEngine(mw ...gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.Use(mw...)
	return r
}

func TestAuth_ValidToken(t *testing.T) {
	validator := TokenValidator(func(ctx context.Context, token string) (*TokenValidateResult, error) {
		return &TokenValidateResult{UserID: 42, JTI: "jti-123"}, nil
	})

	r := newTestEngine(ErrorHandler())
	r.GET("/test", Auth(validator), func(c *gin.Context) {
		uid := c.GetUint("user_id")
		if uid != 42 {
			t.Errorf("expected userID 42, got %d", uid)
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuth_MissingHeader(t *testing.T) {
	validator := TokenValidator(func(ctx context.Context, token string) (*TokenValidateResult, error) {
		return nil, errors.Unauthorized("invalid")
	})

	r := newTestEngine(ErrorHandler())
	r.GET("/test", Auth(validator), func(c *gin.Context) {
		t.Error("handler should not be called")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuth_InvalidFormat(t *testing.T) {
	validator := TokenValidator(func(ctx context.Context, token string) (*TokenValidateResult, error) {
		return nil, nil
	})

	r := newTestEngine(ErrorHandler())
	r.GET("/test", Auth(validator), func(c *gin.Context) {
		t.Error("handler should not be called")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuth_InvalidToken(t *testing.T) {
	validator := TokenValidator(func(ctx context.Context, token string) (*TokenValidateResult, error) {
		return nil, errors.Unauthorized("token expired")
	})

	r := newTestEngine(ErrorHandler())
	r.GET("/test", Auth(validator), func(c *gin.Context) {
		t.Error("handler should not be called")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer expired-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestErrorHandler_CustomError(t *testing.T) {
	r := newTestEngine(ErrorHandler())
	r.GET("/test", func(c *gin.Context) {
		c.Error(errors.NotFound("resource not found"))
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestErrorHandler_GenericError(t *testing.T) {
	r := newTestEngine(ErrorHandler())
	r.GET("/test", func(c *gin.Context) {
		c.Error(context.DeadlineExceeded)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestErrorHandler_NoError(t *testing.T) {
	r := newTestEngine(ErrorHandler())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
