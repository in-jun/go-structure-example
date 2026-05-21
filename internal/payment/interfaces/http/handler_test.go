package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/payment/application/command"
	"github.com/in-jun/go-structure-example/internal/payment/application/query"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/server"
)

type mockCommandUseCase struct {
	createResp *command.CreatePaymentResult
	err        error
}

func (m *mockCommandUseCase) CreatePayment(_ context.Context, _ command.CreatePayment) (*command.CreatePaymentResult, error) {
	return m.createResp, m.err
}
func (m *mockCommandUseCase) ConfirmPayment(_ context.Context, _ command.ConfirmPayment) error {
	return m.err
}
func (m *mockCommandUseCase) RefundPayment(_ context.Context, _ command.RefundPayment) error {
	return m.err
}

type mockQueryUseCase struct {
	getResp *query.Result
	err     error
}

func (m *mockQueryUseCase) GetPayment(_ context.Context, _ query.GetPayment) (*query.Result, error) {
	return m.getResp, m.err
}
func (m *mockQueryUseCase) GetEvents(_ context.Context, _ query.EventHistory) (*query.EventHistoryResult, error) {
	return &query.EventHistoryResult{Events: []query.EventHistoryItem{}}, m.err
}

const testUserID = "550e8400-e29b-41d4-a716-446655440000"
const testPaymentID = "660e8400-e29b-41d4-a716-446655440000"

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

	mux.Handle("GET /api/v1/payments/{id}", noopMw(injectUser(http.HandlerFunc(h.GetPayment))))
	mux.Handle("GET /api/v1/payments/{id}/events", noopMw(injectUser(http.HandlerFunc(h.GetEvents))))
	mux.Handle("POST /api/v1/payments/{id}/confirm", noopMw(injectUser(http.HandlerFunc(h.ConfirmPayment))))
	mux.Handle("POST /api/v1/payments/{id}/refund", noopMw(injectUser(http.HandlerFunc(h.RefundPayment))))

	return mux
}

func TestHandler_GetPayment(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	qryMock := &mockQueryUseCase{
		getResp: &query.Result{
			ID:        testPaymentID,
			AuctionID: "auction-id",
			WinnerID:  testUserID,
			Amount:    5000,
			Status:    "pending",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	router := setupRouter(&mockCommandUseCase{}, qryMock)
	req := httptest.NewRequest("GET", "/api/v1/payments/"+testPaymentID, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestHandler_GetPayment_NotFound(t *testing.T) {
	qryMock := &mockQueryUseCase{
		err: errors.NotFound("Payment not found"),
	}

	router := setupRouter(&mockCommandUseCase{}, qryMock)
	req := httptest.NewRequest("GET", "/api/v1/payments/"+testPaymentID, nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestHandler_ConfirmPayment(t *testing.T) {
	router := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/payments/"+testPaymentID+"/confirm", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_ConfirmPayment_Error(t *testing.T) {
	cmdMock := &mockCommandUseCase{
		err: errors.Forbidden("Not authorized"),
	}

	router := setupRouter(cmdMock, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/payments/"+testPaymentID+"/confirm", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestHandler_RefundPayment(t *testing.T) {
	router := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/payments/"+testPaymentID+"/refund", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_RefundPayment_NotOwner(t *testing.T) {
	cmdMock := &mockCommandUseCase{
		err: errors.Forbidden("Not authorized"),
	}

	router := setupRouter(cmdMock, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/payments/"+testPaymentID+"/refund", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", w.Code)
	}
}

func TestHandler_ConfirmPayment_Conflict(t *testing.T) {
	cmdMock := &mockCommandUseCase{err: errors.Conflict("Payment already processed")}

	router := setupRouter(cmdMock, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/payments/"+testPaymentID+"/confirm", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %d", w.Code)
	}
}

func TestHandler_RefundPayment_Conflict(t *testing.T) {
	cmdMock := &mockCommandUseCase{err: errors.Conflict("Payment cannot be refunded")}

	router := setupRouter(cmdMock, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/payments/"+testPaymentID+"/refund", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected status 409, got %d", w.Code)
	}
}

func TestHandler_GetEvents(t *testing.T) {
	router := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("GET", "/api/v1/payments/"+testPaymentID+"/events", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
