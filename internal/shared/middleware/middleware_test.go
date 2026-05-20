package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
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
	validator := TokenValidator(func(ctx context.Context, token string) (*ValidateTokenResult, error) {
		return &ValidateTokenResult{UserID: 42, JTI: "jti-123"}, nil
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
	validator := TokenValidator(func(ctx context.Context, token string) (*ValidateTokenResult, error) {
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
	validator := TokenValidator(func(ctx context.Context, token string) (*ValidateTokenResult, error) {
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
	validator := TokenValidator(func(ctx context.Context, token string) (*ValidateTokenResult, error) {
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

func TestRequestID_GeneratesID(t *testing.T) {
	r := newTestEngine(RequestID())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Header().Get("X-Request-ID") == "" {
		t.Error("expected X-Request-ID header to be set")
	}
}

func TestRequestID_PropagatesExisting(t *testing.T) {
	r := newTestEngine(RequestID())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", "existing-id")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got := w.Header().Get("X-Request-ID"); got != "existing-id" {
		t.Errorf("expected 'existing-id', got %q", got)
	}
}

func TestSecurityHeaders(t *testing.T) {
	r := newTestEngine(SecurityHeaders())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if v := w.Header().Get("X-Content-Type-Options"); v != "nosniff" {
		t.Errorf("expected 'nosniff', got %q", v)
	}
	if v := w.Header().Get("X-Frame-Options"); v != "DENY" {
		t.Errorf("expected 'DENY', got %q", v)
	}
	if v := w.Header().Get("Referrer-Policy"); v != "strict-origin-when-cross-origin" {
		t.Errorf("expected 'strict-origin-when-cross-origin', got %q", v)
	}
}

func TestCORS_Wildcard(t *testing.T) {
	r := newTestEngine(CORS("*"))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if acao := w.Header().Get("Access-Control-Allow-Origin"); acao != "*" {
		t.Errorf("expected '*', got %q", acao)
	}
}

func TestCORS_SpecificOrigin(t *testing.T) {
	r := newTestEngine(CORS("http://a.com, http://b.com"))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://b.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if acao := w.Header().Get("Access-Control-Allow-Origin"); acao != "http://b.com" {
		t.Errorf("expected 'http://b.com', got %q", acao)
	}
}

func TestCORS_UnknownOrigin(t *testing.T) {
	r := newTestEngine(CORS("http://a.com"))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://evil.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if acao := w.Header().Get("Access-Control-Allow-Origin"); acao != "" {
		t.Errorf("expected empty, got %q", acao)
	}
}

func TestCORS_Preflight(t *testing.T) {
	r := newTestEngine(CORS("*"))
	r.OPTIONS("/test", func(c *gin.Context) {})

	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
}

func TestBodyLimit_Within(t *testing.T) {
	r := newTestEngine(BodyLimit(1024))
	r.POST("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	body := strings.NewReader(`{"name":"test"}`)
	req := httptest.NewRequest("POST", "/test", body)
	req.ContentLength = int64(body.Len())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestBodyLimit_Exceeded(t *testing.T) {
	r := newTestEngine(ErrorHandler(), BodyLimit(10))
	r.POST("/test", func(c *gin.Context) {
		t.Error("handler should not be called")
	})

	body := strings.NewReader(`{"name":"this is way too long"}`)
	req := httptest.NewRequest("POST", "/test", body)
	req.ContentLength = int64(body.Len())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected 413, got %d", w.Code)
	}
}

func TestRateLimit(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}
	defer mr.Close()

	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})
	defer rc.Close()

	r := newTestEngine(RateLimit(rc, 1))
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "1.2.3.4:5678"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("first request: expected 200, got %d", w.Code)
	}

	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req)
	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("second request: expected 429, got %d", w2.Code)
	}
}

func TestTimeout_Normal(t *testing.T) {
	r := newTestEngine(Timeout(1 * time.Second))
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

func TestTimeout_Exceeded(t *testing.T) {
	r := newTestEngine(Timeout(10 * time.Millisecond))
	r.GET("/test", func(c *gin.Context) {
		time.Sleep(100 * time.Millisecond)
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("expected 504, got %d", w.Code)
	}
}

func TestAccessLog(t *testing.T) {
	r := newTestEngine(AccessLog())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "1.2.3.4:5678"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
