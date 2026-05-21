package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, http.StatusCreated, map[string]string{"key": "value"})

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %q", ct)
	}

	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["key"] != "value" {
		t.Errorf("expected 'value', got %q", body["key"])
	}
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, http.StatusBadRequest, "bad input")

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	var body map[string]any
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["message"] != "bad input" {
		t.Errorf("expected 'bad input', got %v", body["message"])
	}
}

func TestBind(t *testing.T) {
	type req struct {
		Name string `json:"name"`
	}
	r := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"test"}`))
	var v req
	if err := Bind(r, &v); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Name != "test" {
		t.Errorf("expected 'test', got %q", v.Name)
	}
}

func TestBind_InvalidJSON(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader(`{invalid}`))
	var v struct{}
	if err := Bind(r, &v); err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestPathParam(t *testing.T) {
	router := NewRouter()
	var got string
	router.HandleFunc("GET /items/{id}", func(w http.ResponseWriter, r *http.Request) {
		got = PathParam(r, "id")
	})

	r := httptest.NewRequest("GET", "/items/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if got != "abc" {
		t.Errorf("expected 'abc', got %q", got)
	}
}

func TestQueryDefault(t *testing.T) {
	r := httptest.NewRequest("GET", "/?page=5", nil)
	if v := QueryDefault(r, "page", "1"); v != "5" {
		t.Errorf("expected '5', got %q", v)
	}
	if v := QueryDefault(r, "limit", "10"); v != "10" {
		t.Errorf("expected default '10', got %q", v)
	}
}

func TestContextWithUserID_And_UserID(t *testing.T) {
	ctx := ContextWithUserID(context.Background(), "user-123")
	r := httptest.NewRequest("GET", "/", nil).WithContext(ctx)
	if id := UserID(r); id != "user-123" {
		t.Errorf("expected 'user-123', got %q", id)
	}
}

func TestUserID_Empty(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	if id := UserID(r); id != "" {
		t.Errorf("expected empty, got %q", id)
	}
}

func TestClientIP(t *testing.T) {
	tests := []struct {
		name   string
		xff    string
		xri    string
		remote string
		want   string
	}{
		{"X-Forwarded-For", "1.2.3.4", "", "5.6.7.8:1234", "1.2.3.4"},
		{"X-Real-Ip", "", "10.0.0.1", "5.6.7.8:1234", "10.0.0.1"},
		{"RemoteAddr", "", "", "5.6.7.8:1234", "5.6.7.8"},
		{"RemoteAddr no port", "", "", "5.6.7.8", "5.6.7.8"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("GET", "/", nil)
			r.RemoteAddr = tt.remote
			if tt.xff != "" {
				r.Header.Set("X-Forwarded-For", tt.xff)
			}
			if tt.xri != "" {
				r.Header.Set("X-Real-Ip", tt.xri)
			}
			if got := ClientIP(r); got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestChain(t *testing.T) {
	var order []string

	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw1-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw1-after")
		})
	}
	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw2-before")
			next.ServeHTTP(w, r)
			order = append(order, "mw2-after")
		})
	}

	handler := Chain(Middleware(mw1), Middleware(mw2))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))

	expected := []string{"mw1-before", "mw2-before", "handler", "mw2-after", "mw1-after"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d: %v", len(expected), len(order), order)
	}
	for i := range expected {
		if order[i] != expected[i] {
			t.Errorf("at %d: expected %q, got %q", i, expected[i], order[i])
		}
	}
}

func TestResponseWriter(t *testing.T) {
	w := httptest.NewRecorder()
	rw := NewResponseWriter(w)

	rw.WriteHeader(http.StatusNotFound)
	_, _ = rw.Write([]byte("not found"))

	if rw.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rw.StatusCode)
	}
	if rw.Body.String() != "not found" {
		t.Errorf("expected 'not found', got %q", rw.Body.String())
	}
}

func TestResponseWriter_DefaultStatus(t *testing.T) {
	w := httptest.NewRecorder()
	rw := NewResponseWriter(w)
	if rw.StatusCode != http.StatusOK {
		t.Errorf("expected default 200, got %d", rw.StatusCode)
	}
}

func TestResponseWriter_DoubleWriteHeader(t *testing.T) {
	w := httptest.NewRecorder()
	rw := NewResponseWriter(w)
	rw.WriteHeader(http.StatusCreated)
	rw.WriteHeader(http.StatusNotFound)
	if rw.StatusCode != http.StatusCreated {
		t.Errorf("expected first status 201, got %d", rw.StatusCode)
	}
}
