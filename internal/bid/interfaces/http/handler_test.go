package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/in-jun/go-structure-example/internal/bid/application/command"
	"github.com/in-jun/go-structure-example/internal/bid/application/query"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/server"
)

type mockCommandUseCase struct {
	placeBidResp *command.PlaceBidResult
	err          error
}

func (m *mockCommandUseCase) PlaceBid(_ context.Context, _ command.PlaceBid) (*command.PlaceBidResult, error) {
	return m.placeBidResp, m.err
}
func (m *mockCommandUseCase) DetermineWinner(_ context.Context, _ command.DetermineWinner) error {
	return m.err
}

type mockQueryUseCase struct {
	highestResp *query.Result
	listResp    *query.ListResult
	err         error
}

func (m *mockQueryUseCase) GetHighest(_ context.Context, _ query.GetHighest) (*query.Result, error) {
	return m.highestResp, m.err
}
func (m *mockQueryUseCase) ListBids(_ context.Context, _ query.ListBids) (*query.ListResult, error) {
	return m.listResp, m.err
}
func (m *mockQueryUseCase) GetEvents(_ context.Context, _ query.EventHistory) (*query.EventHistoryResult, error) {
	return &query.EventHistoryResult{Events: []query.EventHistoryItem{}}, m.err
}

const testUserID = "550e8400-e29b-41d4-a716-446655440000"
const testAuctionID = "660e8400-e29b-41d4-a716-446655440000"

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

	mux.Handle("GET /api/v1/auctions/{auction_id}/bids", noopMw(http.HandlerFunc(h.ListBids)))
	mux.Handle("GET /api/v1/auctions/{auction_id}/bids/highest", noopMw(http.HandlerFunc(h.GetHighest)))
	mux.Handle("GET /api/v1/auctions/{auction_id}/bids/events", noopMw(http.HandlerFunc(h.GetEvents)))
	mux.Handle("POST /api/v1/auctions/{auction_id}/bids", noopMw(injectUser(http.HandlerFunc(h.PlaceBid))))

	return mux
}

func TestHandler_PlaceBid(t *testing.T) {
	cmdMock := &mockCommandUseCase{
		placeBidResp: &command.PlaceBidResult{
			ID:        "bid-id",
			AuctionID: testAuctionID,
			BidderID:  testUserID,
			Amount:    1000,
		},
	}

	router := setupRouter(cmdMock, &mockQueryUseCase{})
	body, _ := json.Marshal(PlaceBidRequest{Amount: 1000})
	req := httptest.NewRequest("POST", "/api/v1/auctions/"+testAuctionID+"/bids", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_PlaceBid_InvalidRequest(t *testing.T) {
	cmdMock := &mockCommandUseCase{
		err: errors.BadRequest("Invalid request"),
	}
	router := setupRouter(cmdMock, &mockQueryUseCase{})
	body, _ := json.Marshal(PlaceBidRequest{Amount: 0})
	req := httptest.NewRequest("POST", "/api/v1/auctions/"+testAuctionID+"/bids", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandler_ListBids(t *testing.T) {
	qryMock := &mockQueryUseCase{
		listResp: &query.ListResult{
			Bids: []query.Result{
				{ID: "b1", Amount: 2000},
				{ID: "b2", Amount: 1500},
			},
			Total: 2,
		},
	}

	router := setupRouter(&mockCommandUseCase{}, qryMock)
	req := httptest.NewRequest("GET", "/api/v1/auctions/"+testAuctionID+"/bids", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp ListResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Total != 2 {
		t.Errorf("expected total 2, got %d", resp.Total)
	}
}

func TestHandler_GetHighest(t *testing.T) {
	qryMock := &mockQueryUseCase{
		highestResp: &query.Result{ID: "b1", Amount: 5000},
	}

	router := setupRouter(&mockCommandUseCase{}, qryMock)
	req := httptest.NewRequest("GET", "/api/v1/auctions/"+testAuctionID+"/bids/highest", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestHandler_GetHighest_NotFound(t *testing.T) {
	qryMock := &mockQueryUseCase{
		err: errors.NotFound("No bids found"),
	}

	router := setupRouter(&mockCommandUseCase{}, qryMock)
	req := httptest.NewRequest("GET", "/api/v1/auctions/"+testAuctionID+"/bids/highest", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestHandler_GetEvents(t *testing.T) {
	router := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("GET", "/api/v1/auctions/"+testAuctionID+"/bids/events", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestHandler_PlaceBid_Forbidden(t *testing.T) {
	cmdMock := &mockCommandUseCase{err: errors.Forbidden("Cannot bid on your own auction")}
	router := setupRouter(cmdMock, &mockQueryUseCase{})
	body, _ := json.Marshal(PlaceBidRequest{Amount: 1000})
	req := httptest.NewRequest("POST", "/api/v1/auctions/"+testAuctionID+"/bids", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestHandler_ListBids_Error(t *testing.T) {
	qryMock := &mockQueryUseCase{err: errors.Internal("db error")}
	router := setupRouter(&mockCommandUseCase{}, qryMock)
	req := httptest.NewRequest("GET", "/api/v1/auctions/"+testAuctionID+"/bids", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}
