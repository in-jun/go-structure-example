package health

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/in-jun/go-structure-example/internal/shared/server"
)

type failingConnector struct{}

func (c *failingConnector) Connect(context.Context) (driver.Conn, error) {
	return nil, errors.New("connection refused")
}

func (c *failingConnector) Driver() driver.Driver { return nil }

func TestLiveHandler(t *testing.T) {
	mux := server.NewRouter()
	NewCheckerWithoutDB(nil).RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/health/live", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %v", body["status"])
	}
}

func TestLiveHandler_WithBuildInfo(t *testing.T) {
	mux := server.NewRouter()
	NewCheckerWithoutDB(nil).WithBuildInfo("v1.0.0", "2026-01-01", "abc123").RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/health/live", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["version"] != "v1.0.0" {
		t.Errorf("expected version v1.0.0, got %v", body["version"])
	}
	if body["git_commit"] != "abc123" {
		t.Errorf("expected git_commit abc123, got %v", body["git_commit"])
	}
}

func TestReadyHandler_NoDeps(t *testing.T) {
	mux := server.NewRouter()
	NewCheckerWithoutDB(nil).RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ready" {
		t.Errorf("expected status ready, got %v", body["status"])
	}
}

func TestReadyHandler_DBDown(t *testing.T) {
	db := sql.OpenDB(&failingConnector{})
	defer db.Close()

	mux := server.NewRouter()
	NewChecker(db, nil).RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d; body: %s", w.Code, w.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "not ready" {
		t.Errorf("expected status 'not ready', got %v", body["status"])
	}
	checks, ok := body["checks"].(map[string]any)
	if !ok {
		t.Fatal("expected checks map in response")
	}
	if _, hasPostgres := checks["postgres"]; !hasPostgres {
		t.Error("expected postgres check in response")
	}
}
