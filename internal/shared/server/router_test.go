package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func dummyHandler(tag string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(tag))
	})
}

func TestRouter_StaticRoutes(t *testing.T) {
	r := NewRouter()
	r.Handle("GET /health/live", dummyHandler("live"))
	r.Handle("GET /health/ready", dummyHandler("ready"))
	r.Handle("GET /metrics", dummyHandler("metrics"))

	tests := []struct {
		path string
		want string
	}{
		{"/health/live", "live"},
		{"/health/ready", "ready"},
		{"/metrics", "metrics"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", tt.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Body.String() != tt.want {
			t.Errorf("GET %s: got %q, want %q", tt.path, w.Body.String(), tt.want)
		}
		if w.Code != 200 {
			t.Errorf("GET %s: got status %d, want 200", tt.path, w.Code)
		}
	}
}

func TestRouter_ParamRoutes(t *testing.T) {
	r := NewRouter()
	r.Handle("GET /api/v1/auctions/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.PathValue("id")))
	}))

	req := httptest.NewRequest("GET", "/api/v1/auctions/abc-123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("got status %d, want 200", w.Code)
	}
	if w.Body.String() != "abc-123" {
		t.Errorf("got %q, want %q", w.Body.String(), "abc-123")
	}
}

func TestRouter_NestedParams(t *testing.T) {
	r := NewRouter()
	r.Handle("GET /api/v1/auctions/{auction_id}/bids", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("auction:" + r.PathValue("auction_id")))
	}))
	r.Handle("GET /api/v1/auctions/{auction_id}/bids/highest", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("highest:" + r.PathValue("auction_id")))
	}))

	tests := []struct {
		path string
		want string
	}{
		{"/api/v1/auctions/a1/bids", "auction:a1"},
		{"/api/v1/auctions/a2/bids/highest", "highest:a2"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", tt.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Body.String() != tt.want {
			t.Errorf("GET %s: got %q, want %q", tt.path, w.Body.String(), tt.want)
		}
	}
}

func TestRouter_MethodRouting(t *testing.T) {
	r := NewRouter()
	r.Handle("GET /api/v1/auctions", dummyHandler("list"))
	r.Handle("POST /api/v1/auctions", dummyHandler("create"))

	req := httptest.NewRequest("GET", "/api/v1/auctions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Body.String() != "list" {
		t.Errorf("GET: got %q, want %q", w.Body.String(), "list")
	}

	req = httptest.NewRequest("POST", "/api/v1/auctions", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Body.String() != "create" {
		t.Errorf("POST: got %q, want %q", w.Body.String(), "create")
	}
}

func TestRouter_NotFound(t *testing.T) {
	r := NewRouter()
	r.Handle("GET /api/v1/auctions", dummyHandler("list"))

	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("got status %d, want 404", w.Code)
	}
}

func TestRouter_MethodNotAllowed(t *testing.T) {
	r := NewRouter()
	r.Handle("GET /api/v1/auctions", dummyHandler("list"))
	r.Handle("POST /api/v1/auctions", dummyHandler("create"))

	req := httptest.NewRequest("DELETE", "/api/v1/auctions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 405 {
		t.Errorf("got status %d, want 405", w.Code)
	}

	allow := w.Header().Get("Allow")
	if allow != "GET, POST" {
		t.Errorf("Allow header: got %q, want %q", allow, "GET, POST")
	}
}

func TestRouter_PathValue(t *testing.T) {
	r := NewRouter()
	var captured string
	r.Handle("GET /users/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.PathValue("id")
	}))

	req := httptest.NewRequest("GET", "/users/user-42", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if captured != "user-42" {
		t.Errorf("PathValue: got %q, want %q", captured, "user-42")
	}
}

func TestRouter_Pattern(t *testing.T) {
	r := NewRouter()
	var captured string
	r.Handle("GET /api/v1/auctions/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = r.Pattern
	}))

	req := httptest.NewRequest("GET", "/api/v1/auctions/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	want := "GET /api/v1/auctions/{id}"
	if captured != want {
		t.Errorf("Pattern: got %q, want %q", captured, want)
	}
}

func TestRouter_AllProjectRoutes(t *testing.T) {
	r := NewRouter()

	routes := []struct {
		pattern string
		tag     string
	}{
		{"GET /health/live", "health-live"},
		{"GET /health/ready", "health-ready"},
		{"GET /metrics", "metrics"},
		{"POST /api/v1/auth/register", "register"},
		{"POST /api/v1/auth/login", "login"},
		{"POST /api/v1/auth/refresh", "refresh"},
		{"POST /api/v1/auth/logout", "logout"},
		{"POST /api/v1/auth/logout/all", "logout-all"},
		{"GET /api/v1/auctions", "auction-list"},
		{"POST /api/v1/auctions", "auction-create"},
		{"GET /api/v1/auctions/{id}", "auction-get"},
		{"GET /api/v1/auctions/{id}/events", "auction-events"},
		{"POST /api/v1/auctions/{id}/open", "auction-open"},
		{"POST /api/v1/auctions/{id}/close", "auction-close"},
		{"GET /api/v1/auctions/{auction_id}/bids", "bid-list"},
		{"GET /api/v1/auctions/{auction_id}/bids/highest", "bid-highest"},
		{"POST /api/v1/auctions/{auction_id}/bids", "bid-place"},
		{"GET /api/v1/payments/{id}", "payment-get"},
		{"POST /api/v1/payments/{id}/confirm", "payment-confirm"},
		{"POST /api/v1/payments/{id}/refund", "payment-refund"},
	}

	for _, rt := range routes {
		r.Handle(rt.pattern, dummyHandler(rt.tag))
	}

	tests := []struct {
		method string
		path   string
		want   string
	}{
		{"GET", "/health/live", "health-live"},
		{"GET", "/health/ready", "health-ready"},
		{"GET", "/metrics", "metrics"},
		{"POST", "/api/v1/auth/register", "register"},
		{"POST", "/api/v1/auth/login", "login"},
		{"POST", "/api/v1/auth/refresh", "refresh"},
		{"POST", "/api/v1/auth/logout", "logout"},
		{"POST", "/api/v1/auth/logout/all", "logout-all"},
		{"GET", "/api/v1/auctions", "auction-list"},
		{"POST", "/api/v1/auctions", "auction-create"},
		{"GET", "/api/v1/auctions/uuid-1", "auction-get"},
		{"GET", "/api/v1/auctions/uuid-1/events", "auction-events"},
		{"POST", "/api/v1/auctions/uuid-1/open", "auction-open"},
		{"POST", "/api/v1/auctions/uuid-1/close", "auction-close"},
		{"GET", "/api/v1/auctions/uuid-1/bids", "bid-list"},
		{"GET", "/api/v1/auctions/uuid-1/bids/highest", "bid-highest"},
		{"POST", "/api/v1/auctions/uuid-1/bids", "bid-place"},
		{"GET", "/api/v1/payments/uuid-2", "payment-get"},
		{"POST", "/api/v1/payments/uuid-2/confirm", "payment-confirm"},
		{"POST", "/api/v1/payments/uuid-2/refund", "payment-refund"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.method, tt.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("%s %s: got status %d, want 200", tt.method, tt.path, w.Code)
		}
		if w.Body.String() != tt.want {
			t.Errorf("%s %s: got %q, want %q", tt.method, tt.path, w.Body.String(), tt.want)
		}
	}
}

func TestRouter_StaticAndParamSamePrefixWithHandler(t *testing.T) {
	r := NewRouter()
	r.Handle("GET /api/v1/auctions", dummyHandler("list"))
	r.Handle("GET /api/v1/auctions/{id}", dummyHandler("get"))
	r.Handle("GET /api/v1/auctions/{id}/events", dummyHandler("events"))

	tests := []struct {
		path string
		want string
	}{
		{"/api/v1/auctions", "list"},
		{"/api/v1/auctions/abc", "get"},
		{"/api/v1/auctions/abc/events", "events"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", tt.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Body.String() != tt.want {
			t.Errorf("GET %s: got %q, want %q", tt.path, w.Body.String(), tt.want)
		}
	}
}

func TestRouter_EmptyParamValue(t *testing.T) {
	r := NewRouter()
	r.Handle("GET /api/{id}", dummyHandler("ok"))

	req := httptest.NewRequest("GET", "/api/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("empty param: got status %d, want 404", w.Code)
	}
}

func TestRouter_CatchAll(t *testing.T) {
	r := NewRouter()
	r.Handle("GET /static/{file...}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.PathValue("file")))
	}))
	r.Handle("GET /api/v1/auctions", dummyHandler("list"))

	tests := []struct {
		path string
		code int
		body string
	}{
		{"/static/css/style.css", 200, "css/style.css"},
		{"/static/js/app.js", 200, "js/app.js"},
		{"/static/favicon.ico", 200, "favicon.ico"},
		{"/static/", 200, ""},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", tt.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tt.code {
			t.Errorf("GET %s: got status %d, want %d", tt.path, w.Code, tt.code)
		}
		if w.Body.String() != tt.body {
			t.Errorf("GET %s: got %q, want %q", tt.path, w.Body.String(), tt.body)
		}
	}
}

func TestRouter_TrailingSlashRedirect(t *testing.T) {
	r := NewRouter()
	r.Handle("GET /api/v1/auctions", dummyHandler("list"))
	r.Handle("GET /health/", dummyHandler("health"))

	// /api/v1/auctions/ → 301 to /api/v1/auctions
	req := httptest.NewRequest("GET", "/api/v1/auctions/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 301 {
		t.Errorf("trailing slash: got %d, want 301", w.Code)
	}

	// /health → 301 to /health/
	req = httptest.NewRequest("GET", "/health", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 301 {
		t.Errorf("missing slash: got %d, want 301", w.Code)
	}
}

func TestRouter_PathClean(t *testing.T) {
	r := NewRouter()
	r.Handle("GET /api/v1/auctions", dummyHandler("list"))

	// /../api/v1/auctions → 301 redirect to /api/v1/auctions
	req := httptest.NewRequest("GET", "/api/v1/../v1/auctions", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 301 {
		t.Errorf("path clean: got %d, want 301", w.Code)
	}

	// //api//v1//auctions → 301 redirect
	req = httptest.NewRequest("GET", "//api//v1//auctions", nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != 301 {
		t.Errorf("double slash: got %d, want 301", w.Code)
	}
}

func TestRouter_HostRouting(t *testing.T) {
	r := NewRouter()
	r.Handle("GET api.example.com/users", dummyHandler("api-users"))
	r.Handle("GET www.example.com/users", dummyHandler("web-users"))
	r.Handle("GET /users", dummyHandler("default-users"))

	tests := []struct {
		host string
		want string
	}{
		{"api.example.com", "api-users"},
		{"www.example.com", "web-users"},
		{"other.com", "default-users"},
	}

	for _, tt := range tests {
		req := httptest.NewRequest("GET", "/users", nil)
		req.Host = tt.host
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Body.String() != tt.want {
			t.Errorf("host %s: got %q, want %q", tt.host, w.Body.String(), tt.want)
		}
	}
}

// Benchmarks

func setupBenchRouter() *Router {
	r := NewRouter()
	r.Handle("GET /health/live", dummyHandler(""))
	r.Handle("GET /health/ready", dummyHandler(""))
	r.Handle("GET /metrics", dummyHandler(""))
	r.Handle("POST /api/v1/auth/register", dummyHandler(""))
	r.Handle("POST /api/v1/auth/login", dummyHandler(""))
	r.Handle("POST /api/v1/auth/refresh", dummyHandler(""))
	r.Handle("POST /api/v1/auth/logout", dummyHandler(""))
	r.Handle("POST /api/v1/auth/logout/all", dummyHandler(""))
	r.Handle("GET /api/v1/auctions", dummyHandler(""))
	r.Handle("POST /api/v1/auctions", dummyHandler(""))
	r.Handle("GET /api/v1/auctions/{id}", dummyHandler(""))
	r.Handle("GET /api/v1/auctions/{id}/events", dummyHandler(""))
	r.Handle("POST /api/v1/auctions/{id}/open", dummyHandler(""))
	r.Handle("POST /api/v1/auctions/{id}/close", dummyHandler(""))
	r.Handle("GET /api/v1/auctions/{auction_id}/bids", dummyHandler(""))
	r.Handle("GET /api/v1/auctions/{auction_id}/bids/highest", dummyHandler(""))
	r.Handle("POST /api/v1/auctions/{auction_id}/bids", dummyHandler(""))
	r.Handle("GET /api/v1/payments/{id}", dummyHandler(""))
	r.Handle("POST /api/v1/payments/{id}/confirm", dummyHandler(""))
	r.Handle("POST /api/v1/payments/{id}/refund", dummyHandler(""))
	r.Handle("GET /api/v1/users/{user_id}/orders/{order_id}", dummyHandler(""))
	return r
}

func setupBenchServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("GET /health/live", dummyHandler(""))
	mux.Handle("GET /health/ready", dummyHandler(""))
	mux.Handle("GET /metrics", dummyHandler(""))
	mux.Handle("POST /api/v1/auth/register", dummyHandler(""))
	mux.Handle("POST /api/v1/auth/login", dummyHandler(""))
	mux.Handle("POST /api/v1/auth/refresh", dummyHandler(""))
	mux.Handle("POST /api/v1/auth/logout", dummyHandler(""))
	mux.Handle("POST /api/v1/auth/logout/all", dummyHandler(""))
	mux.Handle("GET /api/v1/auctions", dummyHandler(""))
	mux.Handle("POST /api/v1/auctions", dummyHandler(""))
	mux.Handle("GET /api/v1/auctions/{id}", dummyHandler(""))
	mux.Handle("GET /api/v1/auctions/{id}/events", dummyHandler(""))
	mux.Handle("POST /api/v1/auctions/{id}/open", dummyHandler(""))
	mux.Handle("POST /api/v1/auctions/{id}/close", dummyHandler(""))
	mux.Handle("GET /api/v1/auctions/{auction_id}/bids", dummyHandler(""))
	mux.Handle("GET /api/v1/auctions/{auction_id}/bids/highest", dummyHandler(""))
	mux.Handle("POST /api/v1/auctions/{auction_id}/bids", dummyHandler(""))
	mux.Handle("GET /api/v1/payments/{id}", dummyHandler(""))
	mux.Handle("POST /api/v1/payments/{id}/confirm", dummyHandler(""))
	mux.Handle("POST /api/v1/payments/{id}/refund", dummyHandler(""))
	mux.Handle("GET /api/v1/users/{user_id}/orders/{order_id}", dummyHandler(""))
	return mux
}

func BenchmarkRouter_Static(b *testing.B) {
	r := setupBenchRouter()
	req := httptest.NewRequest("GET", "/health/live", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkServeMux_Static(b *testing.B) {
	mux := setupBenchServeMux()
	req := httptest.NewRequest("GET", "/health/live", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		mux.ServeHTTP(w, req)
	}
}

func BenchmarkRouter_Param(b *testing.B) {
	r := setupBenchRouter()
	req := httptest.NewRequest("GET", "/api/v1/auctions/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkServeMux_Param(b *testing.B) {
	mux := setupBenchServeMux()
	req := httptest.NewRequest("GET", "/api/v1/auctions/550e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		mux.ServeHTTP(w, req)
	}
}

func BenchmarkRouter_LongPath(b *testing.B) {
	r := setupBenchRouter()
	req := httptest.NewRequest("GET", "/api/v1/auctions/550e8400-e29b-41d4-a716-446655440000/bids/highest", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkServeMux_LongPath(b *testing.B) {
	mux := setupBenchServeMux()
	req := httptest.NewRequest("GET", "/api/v1/auctions/550e8400-e29b-41d4-a716-446655440000/bids/highest", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		mux.ServeHTTP(w, req)
	}
}

func BenchmarkRouter_ShortStatic(b *testing.B) {
	r := setupBenchRouter()
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkServeMux_ShortStatic(b *testing.B) {
	mux := setupBenchServeMux()
	req := httptest.NewRequest("GET", "/metrics", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		mux.ServeHTTP(w, req)
	}
}

func BenchmarkRouter_POST(b *testing.B) {
	r := setupBenchRouter()
	req := httptest.NewRequest("POST", "/api/v1/auctions/550e8400-e29b-41d4-a716-446655440000/bids", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkServeMux_POST(b *testing.B) {
	mux := setupBenchServeMux()
	req := httptest.NewRequest("POST", "/api/v1/auctions/550e8400-e29b-41d4-a716-446655440000/bids", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		mux.ServeHTTP(w, req)
	}
}

func BenchmarkRouter_MultiParam(b *testing.B) {
	r := setupBenchRouter()
	req := httptest.NewRequest("GET", "/api/v1/users/550e8400-e29b-41d4-a716-446655440000/orders/660e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkServeMux_MultiParam(b *testing.B) {
	mux := setupBenchServeMux()
	req := httptest.NewRequest("GET", "/api/v1/users/550e8400-e29b-41d4-a716-446655440000/orders/660e8400-e29b-41d4-a716-446655440000", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		mux.ServeHTTP(w, req)
	}
}

func BenchmarkRouter_Miss404(b *testing.B) {
	r := setupBenchRouter()
	req := httptest.NewRequest("GET", "/api/v1/nonexistent/path", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkServeMux_Miss404(b *testing.B) {
	mux := setupBenchServeMux()
	req := httptest.NewRequest("GET", "/api/v1/nonexistent/path", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		mux.ServeHTTP(w, req)
	}
}

func BenchmarkRouter_Miss405(b *testing.B) {
	r := setupBenchRouter()
	req := httptest.NewRequest("DELETE", "/api/v1/auctions", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkServeMux_Miss405(b *testing.B) {
	mux := setupBenchServeMux()
	req := httptest.NewRequest("DELETE", "/api/v1/auctions", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		mux.ServeHTTP(w, req)
	}
}

func BenchmarkRouter_OverlapPrefix(b *testing.B) {
	r := setupBenchRouter()
	req := httptest.NewRequest("POST", "/api/v1/auth/logout/all", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		r.ServeHTTP(w, req)
	}
}

func BenchmarkServeMux_OverlapPrefix(b *testing.B) {
	mux := setupBenchServeMux()
	req := httptest.NewRequest("POST", "/api/v1/auth/logout/all", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		mux.ServeHTTP(w, req)
	}
}

func BenchmarkTreeWalk_Param(b *testing.B) {
	r := setupBenchRouter()
	root := r.get

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		root.search("/api/v1/auctions/550e8400-e29b-41d4-a716-446655440000", nil)
	}
}

func BenchmarkTreeWalk_LongPath(b *testing.B) {
	r := setupBenchRouter()
	root := r.get

	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		root.search("/api/v1/auctions/550e8400-e29b-41d4-a716-446655440000/bids/highest", nil)
	}
}
