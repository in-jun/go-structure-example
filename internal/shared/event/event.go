package event

import "time"

type Event interface {
	EventName() string
	AggregateID() string
	OccurredAt() time.Time
}

type StoredEvent struct {
	ID          int64
	AggregateID string
	EventType   string
	Payload     []byte
	OccurredAt  time.Time
}
