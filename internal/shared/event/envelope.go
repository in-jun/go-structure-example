package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Envelope struct {
	ID           string            `json:"id"`
	Type         string            `json:"type"`
	AggregateID  string            `json:"aggregate_id"`
	Payload      json.RawMessage   `json:"payload"`
	OccurredAt   time.Time         `json:"occurred_at"`
	TraceContext map[string]string `json:"trace_context,omitempty"`
}

func NewEnvelope(eventType, aggregateID string, payload any, occurredAt time.Time) (*Envelope, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &Envelope{
		ID:          uuid.New().String(),
		Type:        eventType,
		AggregateID: aggregateID,
		Payload:     data,
		OccurredAt:  occurredAt,
	}, nil
}

func NewEnvelopeWithID(id, eventType, aggregateID string, payload any, occurredAt time.Time) (*Envelope, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return &Envelope{
		ID:          id,
		Type:        eventType,
		AggregateID: aggregateID,
		Payload:     data,
		OccurredAt:  occurredAt,
	}, nil
}
