package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/in-jun/go-structure-example/internal/auth/application"
	"github.com/in-jun/go-structure-example/internal/auth/application/command"
	"github.com/in-jun/go-structure-example/internal/auth/application/query"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/middleware"
	"github.com/in-jun/go-structure-example/internal/shared/server"
)

const testUUID = "550e8400-e29b-41d4-a716-446655440000"

type mockCommandUseCase struct {
	loginResp   *command.LoginResult
	refreshResp *command.RefreshResult
	err         error
}

func (m *mockCommandUseCase) Register(_ context.Context, _ command.Register) error { return m.err }
func (m *mockCommandUseCase) Login(_ context.Context, _ command.Login) (*command.LoginResult, error) {
	return m.loginResp, m.err
}
func (m *mockCommandUseCase) Refresh(_ context.Context, _ command.Refresh) (*command.RefreshResult, error) {
	return m.refreshResp, m.err
}
func (m *mockCommandUseCase) Logout(_ context.Context, _ command.Logout) error       { return m.err }
func (m *mockCommandUseCase) LogoutAll(_ context.Context, _ command.LogoutAll) error { return m.err }

type mockQueryUseCase struct {
	validateResp *query.Result
	err          error
}

func (m *mockQueryUseCase) ValidateToken(_ context.Context, _ query.Validate) (*query.Result, error) {
	return m.validateResp, m.err
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

func TestHandler_Register(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	body, _ := json.Marshal(RegisterRequest{Email: "test@example.com", Password: "password123", Name: "Test"})
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Register_Error(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{err: errors.Conflict("email exists")}, &mockQueryUseCase{})
	body, _ := json.Marshal(RegisterRequest{Email: "test@example.com", Password: "password123", Name: "Test"})
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestHandler_Login(t *testing.T) {
	cmdMock := &mockCommandUseCase{
		loginResp: &command.LoginResult{
			AccessToken:  "access",
			RefreshToken: "refresh",
			ExpiresIn:    3600,
		},
	}
	mux := setupRouter(cmdMock, &mockQueryUseCase{})
	body, _ := json.Marshal(LoginRequest{Email: "test@example.com", Password: "password"})
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.AccessToken != "access" {
		t.Errorf("expected access token 'access', got %q", resp.AccessToken)
	}
}

func TestHandler_Login_InvalidCredentials(t *testing.T) {
	cmdMock := &mockCommandUseCase{err: errors.Unauthorized("invalid credentials")}
	mux := setupRouter(cmdMock, &mockQueryUseCase{})
	body, _ := json.Marshal(LoginRequest{Email: "test@example.com", Password: "wrong"})
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestHandler_Refresh(t *testing.T) {
	cmdMock := &mockCommandUseCase{
		refreshResp: &command.RefreshResult{
			AccessToken:  "new-access",
			RefreshToken: "new-refresh",
			ExpiresIn:    3600,
		},
	}
	mux := setupRouter(cmdMock, &mockQueryUseCase{})
	body, _ := json.Marshal(RefreshRequest{RefreshToken: "old-refresh"})
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Logout(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	body, _ := json.Marshal(LogoutRequest{RefreshToken: "token"})
	req := httptest.NewRequest("POST", "/api/v1/auth/logout", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_LogoutAll(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/auth/logout/all", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
}

func TestHandler_Register_BadJSON(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader([]byte("{invalid}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandler_Login_BadJSON(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader([]byte("{invalid}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandler_Refresh_BadJSON(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewReader([]byte("{invalid}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandler_Logout_BadJSON(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/auth/logout", bytes.NewReader([]byte("{invalid}")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer valid-token")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestHandler_Logout_MissingAuth(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	body, _ := json.Marshal(LogoutRequest{RefreshToken: "token"})
	req := httptest.NewRequest("POST", "/api/v1/auth/logout", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestHandler_LogoutAll_MissingAuth(t *testing.T) {
	mux := setupRouter(&mockCommandUseCase{}, &mockQueryUseCase{})
	req := httptest.NewRequest("POST", "/api/v1/auth/logout/all", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}
