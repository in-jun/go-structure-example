package event

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewEnvelope(t *testing.T) {
	payload := map[string]string{"key": "value"}
	now := time.Now()

	env, err := NewEnvelope("test.event", "agg-123", payload, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.ID == "" {
		t.Error("expected non-empty ID")
	}
	if env.Type != "test.event" {
		t.Errorf("expected type 'test.event', got %q", env.Type)
	}
	if env.AggregateID != "agg-123" {
		t.Errorf("expected aggregate ID 'agg-123', got %q", env.AggregateID)
	}

	var decoded map[string]string
	if err := json.Unmarshal(env.Payload, &decoded); err != nil {
		t.Fatalf("failed to unmarshal payload: %v", err)
	}
	if decoded["key"] != "value" {
		t.Errorf("expected payload key 'value', got %q", decoded["key"])
	}
}

func TestNewEnvelopeWithID(t *testing.T) {
	payload := map[string]int{"amount": 1000}
	now := time.Now()

	env, err := NewEnvelopeWithID("custom-id", "test.event", "agg-456", payload, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.ID != "custom-id" {
		t.Errorf("expected ID 'custom-id', got %q", env.ID)
	}
}

func TestNewEnvelope_InvalidPayload(t *testing.T) {
	_, err := NewEnvelope("test.event", "agg", make(chan int), time.Now())
	if err == nil {
		t.Error("expected error for unmarshalable payload")
	}
}

func TestNewEnvelopeWithID_InvalidPayload(t *testing.T) {
	_, err := NewEnvelopeWithID("id", "test.event", "agg", make(chan int), time.Now())
	if err == nil {
		t.Error("expected error for unmarshalable payload")
	}
}
