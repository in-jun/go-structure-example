package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/server"
)

func TestRequestID_GeneratesID(t *testing.T) {
	handler := RequestID()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if id := w.Header().Get("X-Request-ID"); id == "" {
		t.Error("expected X-Request-ID header")
	}
}

func TestRequestID_PreservesExisting(t *testing.T) {
	handler := RequestID()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Request-ID", "existing-id")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if id := w.Header().Get("X-Request-ID"); id != "existing-id" {
		t.Errorf("expected 'existing-id', got %q", id)
	}
}

func TestSecurityHeaders(t *testing.T) {
	handler := SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	expected := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Referrer-Policy":        "strict-origin-when-cross-origin",
		"Permissions-Policy":     "camera=(), microphone=(), geolocation=()",
	}
	for k, v := range expected {
		if got := w.Header().Get(k); got != v {
			t.Errorf("%s: expected %q, got %q", k, v, got)
		}
	}
}

func TestCORS_Preflight(t *testing.T) {
	handler := CORS("*")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest("OPTIONS", "/", nil)
	r.Header.Set("Origin", "http://example.com")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d", w.Code)
	}
	if acao := w.Header().Get("Access-Control-Allow-Origin"); acao != "*" {
		t.Errorf("expected '*', got %q", acao)
	}
}

func TestCORS_SpecificOrigin(t *testing.T) {
	handler := CORS("http://a.com, http://b.com")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Origin", "http://b.com")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if acao := w.Header().Get("Access-Control-Allow-Origin"); acao != "http://b.com" {
		t.Errorf("expected 'http://b.com', got %q", acao)
	}
}

func TestCORS_UnknownOrigin(t *testing.T) {
	handler := CORS("http://a.com")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Origin", "http://evil.com")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if acao := w.Header().Get("Access-Control-Allow-Origin"); acao != "" {
		t.Errorf("expected empty, got %q", acao)
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

	handler := RateLimit(rc, 1, 1)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "1.2.3.4:5678"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("first request: expected 200, got %d", w.Code)
	}

	w2 := httptest.NewRecorder()
	handler.ServeHTTP(w2, r)
	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("second request: expected 429, got %d", w2.Code)
	}
}

func TestRecovery(t *testing.T) {
	handler := Recovery()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestAuth_ValidToken(t *testing.T) {
	validator := TokenValidator(func(ctx context.Context, token string) (*ValidateTokenResult, error) {
		return &ValidateTokenResult{UserID: "user-123"}, nil
	})

	handler := Auth(validator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := server.UserID(r)
		if uid != "user-123" {
			t.Errorf("expected user-123, got %q", uid)
		}
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuth_MissingHeader(t *testing.T) {
	validator := TokenValidator(func(ctx context.Context, token string) (*ValidateTokenResult, error) {
		return nil, errors.Unauthorized("invalid")
	})

	handler := Auth(validator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuth_InvalidFormat(t *testing.T) {
	validator := TokenValidator(func(ctx context.Context, token string) (*ValidateTokenResult, error) {
		return nil, nil
	})

	handler := Auth(validator)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "InvalidFormat")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGatewayAuth_Valid(t *testing.T) {
	handler := GatewayAuth()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := server.UserID(r)
		if uid != "550e8400-e29b-41d4-a716-446655440000" {
			t.Errorf("expected valid uuid, got %q", uid)
		}
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-User-ID", "550e8400-e29b-41d4-a716-446655440000")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGatewayAuth_Missing(t *testing.T) {
	handler := GatewayAuth()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not be called")
	}))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGatewayAuth_InvalidUUID(t *testing.T) {
	handler := GatewayAuth()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("should not be called")
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-User-ID", "not-a-uuid")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestHandleError_CustomError(t *testing.T) {
	w := httptest.NewRecorder()
	HandleError(w, errors.NotFound("not found"))

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}

	var body map[string]any
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["message"] != "not found" {
		t.Errorf("expected 'not found', got %v", body["message"])
	}
}

func TestHandleError_GenericError(t *testing.T) {
	w := httptest.NewRecorder()
	HandleError(w, context.DeadlineExceeded)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestAccessLog(t *testing.T) {
	handler := AccessLog()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest("GET", "/test", nil)
	r.RemoteAddr = "1.2.3.4:5678"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestMetrics(t *testing.T) {
	handler := Metrics("test-svc")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestTimeout_Normal(t *testing.T) {
	handler := Timeout(1 * time.Second)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestTimeout_Exceeded(t *testing.T) {
	handler := Timeout(10 * time.Millisecond)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", w.Code)
	}
}

func TestBodyLimit_Within(t *testing.T) {
	handler := BodyLimit(1024)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))

	body := strings.NewReader(`{"name":"test"}`)
	r := httptest.NewRequest("POST", "/", body)
	r.ContentLength = int64(body.Len())
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestBodyLimit_Exceeded(t *testing.T) {
	handler := BodyLimit(10)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called")
	}))

	body := strings.NewReader(`{"name":"this is way too long"}`)
	r := httptest.NewRequest("POST", "/", body)
	r.ContentLength = int64(body.Len())
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("expected 413, got %d", w.Code)
	}
}

func TestRetryTransport_NoRetryOnPOST(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	rt := NewRetryTransport(ts.Client().Transport, 3)
	req, _ := http.NewRequest("POST", ts.URL+"/create", nil)
	resp, _ := rt.RoundTrip(req)
	resp.Body.Close()

	if calls != 1 {
		t.Errorf("POST should not retry, expected 1 call, got %d", calls)
	}
}

func TestRetryTransport_RetryOnGET(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 3 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	rt := NewRetryTransport(ts.Client().Transport, 3)
	req, _ := http.NewRequest("GET", ts.URL+"/data", nil)
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 after retry, got %d", resp.StatusCode)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestRetryTransport_POSTWithIdempotencyKey(t *testing.T) {
	calls := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if calls < 2 {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		fmt.Fprintln(w, "ok")
	}))
	defer ts.Close()

	rt := NewRetryTransport(ts.Client().Transport, 3)
	req, _ := http.NewRequest("POST", ts.URL+"/create", nil)
	req.Header.Set("Idempotency-Key", "key-123")
	resp, err := rt.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}
