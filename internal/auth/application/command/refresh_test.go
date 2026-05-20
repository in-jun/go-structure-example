package command

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/auth/domain/entity"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

func makeValidRefreshToken() *entity.RefreshToken {
	rt, _ := entity.NewRefreshToken(testUUID, time.Now().Add(time.Hour))
	return rt
}

func makeExpiredRefreshToken() *entity.RefreshToken {
	rt, _ := entity.ReconstructRefreshToken("expired-token", testUUID, time.Now().Add(-time.Hour))
	return rt
}

func TestRefreshHandler_Success(t *testing.T) {
	rt := makeValidRefreshToken()
	h := NewRefreshHandler(&mockTokenRepo{token: rt}, &mockTokenGen{})

	result, err := h.Handle(context.Background(), Refresh{RefreshToken: rt.Token()})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if result.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if result.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
}

func TestRefreshHandler_InvalidToken(t *testing.T) {
	h := NewRefreshHandler(&mockTokenRepo{}, &mockTokenGen{})

	_, err := h.Handle(context.Background(), Refresh{RefreshToken: ""})
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestRefreshHandler_TokenNotFound(t *testing.T) {
	h := NewRefreshHandler(&mockTokenRepo{token: nil}, &mockTokenGen{})

	_, err := h.Handle(context.Background(), Refresh{RefreshToken: "nonexistent-token"})
	if err == nil {
		t.Fatal("expected error for not-found token, got nil")
	}
	var ce errors.CustomError
	if !asCustomError(err, &ce) || ce.Status != 401 {
		t.Errorf("expected 401 Unauthorized, got %v", err)
	}
}

func TestRefreshHandler_ExpiredToken(t *testing.T) {
	rt := makeExpiredRefreshToken()
	h := NewRefreshHandler(&mockTokenRepo{token: rt}, &mockTokenGen{})

	_, err := h.Handle(context.Background(), Refresh{RefreshToken: rt.Token()})
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
	var ce errors.CustomError
	if !asCustomError(err, &ce) || ce.Status != 401 {
		t.Errorf("expected 401 Unauthorized, got %v", err)
	}
}

func TestRefreshHandler_TokenReuseDetected(t *testing.T) {
	h := NewRefreshHandler(&mockTokenRepo{token: nil, userID: testUUID}, &mockTokenGen{})

	_, err := h.Handle(context.Background(), Refresh{RefreshToken: "used-token"})
	if err == nil {
		t.Fatal("expected error for reused token, got nil")
	}
	var ce errors.CustomError
	if !asCustomError(err, &ce) || ce.Status != 401 {
		t.Errorf("expected 401 Unauthorized, got %v", err)
	}
}

func TestRefreshHandler_RepositoryError(t *testing.T) {
	h := NewRefreshHandler(&mockTokenRepo{err: errors.Internal("db error")}, &mockTokenGen{})

	_, err := h.Handle(context.Background(), Refresh{RefreshToken: "some-token"})
	if err == nil {
		t.Fatal("expected error for repository failure, got nil")
	}
}
