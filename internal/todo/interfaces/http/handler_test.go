package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/middleware"
	"github.com/in-jun/go-structure-example/internal/shared/server"
	"github.com/in-jun/go-structure-example/internal/todo/application"
	"github.com/in-jun/go-structure-example/internal/todo/application/command"
	"github.com/in-jun/go-structure-example/internal/todo/application/query"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

const testUUID = "550e8400-e29b-41d4-a716-446655440000"

type mockCommandUseCase struct {
	createResult *command.CreateResult
	err          error
}

func (m *mockCommandUseCase) Create(_ context.Context, _ command.Create) (*command.CreateResult, error) {
	return m.createResult, m.err
}
func (m *mockCommandUseCase) Update(_ context.Context, _ command.Update) error { return m.err }
func (m *mockCommandUseCase) UpdateStatus(_ context.Context, _ command.UpdateStatus) error {
	return m.err
}
func (m *mockCommandUseCase) Delete(_ context.Context, _ command.Delete) error { return m.err }

type mockQueryUseCase struct {
	todoResult     *query.Result
	todoListResult *query.ListResult
	err            error
}

func (m *mockQueryUseCase) Get(_ context.Context, _ query.Get) (*query.Result, error) {
	return m.todoResult, m.err
}
func (m *mockQueryUseCase) GetList(_ context.Context, _ query.List) (*query.ListResult, error) {
	return m.todoListResult, m.err
}

var _ application.CommandUseCase = (*mockCommandUseCase)(nil)
var _ application.QueryUseCase = (*mockQueryUseCase)(nil)

func setupRouter(cmdMock *mockCommandUseCase, qryMock *mockQueryUseCase) *server.Router {
	tokenValidator := middleware.TokenValidator(func(ctx context.Context, token string) (*middleware.ValidateTokenResult, error) {
		return &middleware.ValidateTokenResult{UserID: testUUID, JTI: "test-jti"}, nil
	})

	h := NewHandler(cmdMock, qryMock, tokenValidator)
	mux := server.NewRouter()
	h.RegisterRoutes(mux, server.Chain())
	return mux
}

func TestHandler_Create(t *testing.T) {
	cmdMock := &mockCommandUseCase{createResult: &command.CreateResult{ID: testUUID}}
	mux := setupRouter(cmdMock, &mockQueryUseCase{})
	body, _ := json.Marshal(CreateTodoRequest{
		Title:   "Buy groceries",
		DueDate: time.Now().Add(time.Hour),
	})
	req := httptest.NewRequest("POST", "/api/v1/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Create_Error(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{err: errors.BadRequest("title required")}, &mockQueryUseCase{})
	body, _ := json.Marshal(CreateTodoRequest{Title: "", DueDate: time.Now().Add(time.Hour)})
	req := httptest.NewRequest("POST", "/api/v1/todos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandler_GetList(t *testing.T) {
	qryMock := &mockQueryUseCase{
		todoListResult: &query.ListResult{
			Todos: []query.Result{{
				ID:      testUUID,
				Title:   "Test",
				DueDate: time.Now().Add(time.Hour),
			}},
			Total: 1,
		},
	}
	mux := setupRouter(&mockCommandUseCase{}, qryMock)
	req := httptest.NewRequest("GET", "/api/v1/todos", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp TodoListResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Total != 1 {
		t.Errorf("expected total 1, got %d", resp.Total)
	}
}

func TestHandler_Get(t *testing.T) {
	qryMock := &mockQueryUseCase{
		todoResult: &query.Result{
			ID:      testUUID,
			Title:   "Test",
			Status:  entity.StatusPending,
			DueDate: time.Now().Add(time.Hour),
		},
	}
	mux := setupRouter(&mockCommandUseCase{}, qryMock)
	req := httptest.NewRequest("GET", "/api/v1/todos/"+testUUID, nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Get_NotFound(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{err: errors.NotFound("not found")})
	req := httptest.NewRequest("GET", "/api/v1/todos/"+testUUID, nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandler_Update(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	body, _ := json.Marshal(UpdateTodoRequest{Title: "Updated", DueDate: time.Now().Add(time.Hour)})
	req := httptest.NewRequest("PUT", "/api/v1/todos/"+testUUID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdateStatus(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	body, _ := json.Marshal(UpdateTodoStatusRequest{Status: entity.StatusCompleted})
	req := httptest.NewRequest("PATCH", "/api/v1/todos/"+testUUID+"/status", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Delete(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("DELETE", "/api/v1/todos/"+testUUID, nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected 204, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Create_BadJSON(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/todos", bytes.NewReader([]byte("{invalid}")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandler_Update_BadJSON(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("PUT", "/api/v1/todos/"+testUUID, bytes.NewReader([]byte("{invalid}")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandler_UpdateStatus_BadJSON(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("PATCH", "/api/v1/todos/"+testUUID+"/status", bytes.NewReader([]byte("{invalid}")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandler_Get_Forbidden(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{err: errors.Forbidden("access denied")})
	req := httptest.NewRequest("GET", "/api/v1/todos/"+testUUID, nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestHandler_MissingAuth(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("GET", "/api/v1/todos", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
