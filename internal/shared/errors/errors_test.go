package errors

import (
	"net/http"
	"testing"
)

func TestCustomError_Error(t *testing.T) {
	err := BadRequest("bad")
	want := "status: 400, code: BAD_REQUEST, message: bad"
	if err.Error() != want {
		t.Errorf("unexpected: %s", err.Error())
	}
}

func TestErrorFactories(t *testing.T) {
	tests := []struct {
		name       string
		fn         func(string) CustomError
		wantStatus int
		wantCode   string
	}{
		{"BadRequest", BadRequest, http.StatusBadRequest, "BAD_REQUEST"},
		{"Unauthorized", Unauthorized, http.StatusUnauthorized, "UNAUTHORIZED"},
		{"Forbidden", Forbidden, http.StatusForbidden, "FORBIDDEN"},
		{"NotFound", NotFound, http.StatusNotFound, "NOT_FOUND"},
		{"Conflict", Conflict, http.StatusConflict, "CONFLICT"},
		{"TooManyRequests", TooManyRequests, http.StatusTooManyRequests, "TOO_MANY_REQUESTS"},
		{"Internal", Internal, http.StatusInternalServerError, "INTERNAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn("test message")
			if err.Status != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, err.Status)
			}
			if err.Code != tt.wantCode {
				t.Errorf("expected code %q, got %q", tt.wantCode, err.Code)
			}
			if err.Message != "test message" {
				t.Errorf("expected message 'test message', got %q", err.Message)
			}
		})
	}
}

func TestCustomError_Is(t *testing.T) {
	err := NotFound("not found")
	target := NotFound("different message")
	if !err.Is(target) {
		t.Error("expected Is to return true for same status")
	}

	other := BadRequest("bad")
	if err.Is(other) {
		t.Error("expected Is to return false for different status")
	}
}
