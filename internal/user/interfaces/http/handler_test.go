package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/middleware"
	"github.com/in-jun/go-structure-example/internal/user/application/command"
	"github.com/in-jun/go-structure-example/internal/user/application/query"
)

type mockCommandUseCase struct {
	err error
}

func (m *mockCommandUseCase) UpdateProfile(_ context.Context, _ command.UpdateProfile) error {
	return m.err
}
func (m *mockCommandUseCase) UpdatePassword(_ context.Context, _ command.UpdatePassword) error {
	return m.err
}
func (m *mockCommandUseCase) Delete(_ context.Context, _ command.Delete) error { return m.err }

type mockQueryUseCase struct {
	userResult *query.Result
	err        error
}

func (m *mockQueryUseCase) GetProfile(_ context.Context, _ query.Get) (*query.Result, error) {
	return m.userResult, m.err
}

func setupRouter(cmdMock *mockCommandUseCase, qryMock *mockQueryUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())

	tokenValidator := middleware.TokenValidator(func(ctx context.Context, token string) (*middleware.ValidateTokenResult, error) {
		return &middleware.ValidateTokenResult{UserID: 1, JTI: "test-jti"}, nil
	})

	h := NewHandler(cmdMock, qryMock, tokenValidator)
	api := r.Group("/api/v1")
	h.RegisterRoutes(api)
	return r
}

func TestHandler_GetMe(t *testing.T) {
	qryMock := &mockQueryUseCase{
		userResult: &query.Result{
			ID:        1,
			Email:     "test@example.com",
			Name:      "Test User",
			CreatedAt: time.Now(),
		},
	}
	r := setupRouter(&mockCommandUseCase{}, qryMock)
	req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp UserResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Email != "test@example.com" {
		t.Errorf("expected email test@example.com, got %q", resp.Email)
	}
}

func TestHandler_GetMe_NotFound(t *testing.T) {
	r := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{err: errors.NotFound("user not found")})
	req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestHandler_UpdateProfile(t *testing.T) {
	r := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	body, _ := json.Marshal(UpdateProfileRequest{Name: "New Name"})
	req := httptest.NewRequest("PATCH", "/api/v1/users/me/profile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdatePassword(t *testing.T) {
	r := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	body, _ := json.Marshal(UpdatePasswordRequest{CurrentPassword: "old", NewPassword: "newpass"})
	req := httptest.NewRequest("PATCH", "/api/v1/users/me/password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_UpdatePassword_Error(t *testing.T) {
	r := setupRouter(&mockCommandUseCase{err: errors.Unauthorized("wrong password")}, &mockQueryUseCase{})
	body, _ := json.Marshal(UpdatePasswordRequest{CurrentPassword: "wrong", NewPassword: "newpass"})
	req := httptest.NewRequest("PATCH", "/api/v1/users/me/password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestHandler_DeleteMe(t *testing.T) {
	r := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("DELETE", "/api/v1/users/me", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_MissingAuth(t *testing.T) {
	r := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("GET", "/api/v1/users/me", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
