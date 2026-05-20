package logging

import (
	"context"
	"testing"
)

func TestWithRequestID_RoundTrip(t *testing.T) {
	ctx := context.Background()
	ctx = WithRequestID(ctx, "test-request-id")
	if got := RequestIDFromContext(ctx); got != "test-request-id" {
		t.Errorf("expected 'test-request-id', got %q", got)
	}
}

func TestRequestIDFromContext_Empty(t *testing.T) {
	if got := RequestIDFromContext(context.Background()); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestFromContext_WithRequestID(t *testing.T) {
	ctx := WithRequestID(context.Background(), "req-123")
	if logger := FromContext(ctx); logger == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestFromContext_WithoutRequestID(t *testing.T) {
	if logger := FromContext(context.Background()); logger == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestInit(t *testing.T) {
	Init("test-service")
	if logger := FromContext(context.Background()); logger == nil {
		t.Fatal("expected non-nil logger after Init")
	}
}
