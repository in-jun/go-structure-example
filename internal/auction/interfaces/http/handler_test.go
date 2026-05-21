package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/auction/application/command"
	"github.com/in-jun/go-structure-example/internal/auction/application/query"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/server"
)

type mockCommandUseCase struct {
	createResp *command.CreateResult
	err        error
}

func (m *mockCommandUseCase) Create(_ context.Context, _ command.Create) (*command.CreateResult, error) {
	return m.createResp, m.err
}
func (m *mockCommandUseCase) Open(_ context.Context, _ command.Open) error     { return m.err }
func (m *mockCommandUseCase) Close(_ context.Context, _ command.Close) error   { return m.err }
func (m *mockCommandUseCase) Settle(_ context.Context, _ command.Settle) error { return m.err }
func (m *mockCommandUseCase) Cancel(_ context.Context, _ command.Cancel) error { return m.err }

type mockQueryUseCase struct {
	getResp    *query.Result
	listResp   *query.ListResult
	eventsResp *query.EventHistoryResult
	err        error
}

func (m *mockQueryUseCase) GetByID(_ context.Context, _ query.Get) (*query.Result, error) {
	return m.getResp, m.err
}
func (m *mockQueryUseCase) GetList(_ context.Context, _ query.List) (*query.ListResult, error) {
	return m.listResp, m.err
}
func (m *mockQueryUseCase) GetEvents(_ context.Context, _ query.EventHistory) (*query.EventHistoryResult, error) {
	if m.eventsResp != nil {
		return m.eventsResp, m.err
	}
	return &query.EventHistoryResult{Events: []query.EventHistoryItem{}}, m.err
}

const testUserID = "550e8400-e29b-41d4-a716-446655440000"

func setupRouter(cmdMock *mockCommandUseCase, qryMock *mockQueryUseCase) *server.Router {
	mux := server.NewRouter()
	h := NewHandler(cmdMock, qryMock)
	noopMw := server.Middleware(func(next http.Handler) http.Handler { return next })

	injectUser := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := server.ContextWithUserID(r.Context(), testUserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	mux.Handle("GET /api/v1/auctions", noopMw(http.HandlerFunc(h.GetList)))
	mux.Handle("GET /api/v1/auctions/{id}", noopMw(http.HandlerFunc(h.GetByID)))
	mux.Handle("GET /api/v1/auctions/{id}/events", noopMw(http.HandlerFunc(h.GetEvents)))
	mux.Handle("POST /api/v1/auctions", noopMw(injectUser(http.HandlerFunc(h.Create))))
	mux.Handle("POST /api/v1/auctions/{id}/open", noopMw(injectUser(http.HandlerFunc(h.Open))))
	mux.Handle("POST /api/v1/auctions/{id}/close", noopMw(injectUser(http.HandlerFunc(h.Close))))
	mux.Handle("POST /api/v1/auctions/{id}/cancel", noopMw(http.HandlerFunc(h.Cancel)))

	return mux
}

func TestHandler_Create(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	cmdMock := &mockCommandUseCase{
		createResp: &command.CreateResult{
			ID:         "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			Title:      "Test Auction",
			StartPrice: 1000,
			Status:     "draft",
			EndTime:    now.Add(24 * time.Hour),
			CreatedAt:  now,
			UpdatedAt:  now,
		},
	}

	router := setupRouter(cmdMock, &mockQueryUseCase{})
	body, _ := json.Marshal(CreateRequest{
		Title:      "Test Auction",
		StartPrice: 1000,
		EndTime:    now.Add(24 * time.Hour),
	})
	req := httptest.NewRequest("POST", "/api/v1/auctions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Title != "Test Auction" {
		t.Errorf("expected title 'Test Auction', got '%s'", resp.Title)
	}
}

func TestHandler_Create_InvalidRequest(t *testing.T) {
	cmdMock := &mockCommandUseCase{
		err: errors.BadRequest("Invalid request"),
	}
	router := setupRouter(cmdMock, &mockQueryUseCase{})
	body, _ := json.Marshal(CreateRequest{Title: ""})
	req := httptest.NewRequest("POST", "/api/v1/auctions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandler_GetList(t *testing.T) {
	qryMock := &mockQueryUseCase{
		listResp: &query.ListResult{
			Auctions: []query.Result{
				{ID: "6ba7b810-9dad-11d1-80b4-00c04fd430c8", Title: "Auction 1", Status: "open"},
				{ID: "6ba7b811-9dad-11d1-80b4-00c04fd430c8", Title: "Auction 2", Status: "draft"},
			},
			Total: 2,
		},
	}

	router := setupRouter(&mockCommandUseCase{}, qryMock)
	req := httptest.NewRequest("GET", "/api/v1/auctions?page=1&limit=10", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp ListResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

func TestHandler_GetByID(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	qryMock := &mockQueryUseCase{
		getResp: &query.Result{
			ID:       "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			SellerID: testUserID,
			Title:    "Test",
			Status:   "open",
			EndTime:  now.Add(24 * time.Hour),
		},
	}

	router := setupRouter(&mockCommandUseCase{}, qryMock)
	req := httptest.NewRequest("GET", "/api/v1/auctions/6ba7b810-9dad-11d1-80b4-00c04fd430c8", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestHandler_GetByID_NotFound(t *testing.T) {
	qryMock := &mockQueryUseCase{
		err: errors.NotFound("Auction not found"),
	}

	router := setupRouter(&mockCommandUseCase{}, qryMock)
	req := httptest.NewRequest("GET", "/api/v1/auctions/6ba7b810-9dad-11d1-80b4-00c04fd430c8", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestHandler_Open(t *testing.T) {
	router := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/auctions/6ba7b810-9dad-11d1-80b4-00c04fd430c8/open", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Close(t *testing.T) {
	router := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/auctions/6ba7b810-9dad-11d1-80b4-00c04fd430c8/close", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_GetEvents(t *testing.T) {
	router := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("GET", "/api/v1/auctions/6ba7b810-9dad-11d1-80b4-00c04fd430c8/events", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestHandler_Cancel(t *testing.T) {
	router := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/auctions/6ba7b810-9dad-11d1-80b4-00c04fd430c8/cancel", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d; body: %s", w.Code, w.Body.String())
	}
}
