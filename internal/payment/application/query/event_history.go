package query

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/payment/domain"
	"github.com/in-jun/go-structure-example/internal/payment/domain/vo"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

type EventHistory struct {
	PaymentID string
}

type EventHistoryItem struct {
	ID         int64
	EventType  string
	Payload    []byte
	OccurredAt time.Time
}

type EventHistoryResult struct {
	Events []EventHistoryItem
}

type EventHistoryHandler struct {
	eventReader domain.EventReader
}

func NewEventHistoryHandler(eventReader domain.EventReader) *EventHistoryHandler {
	return &EventHistoryHandler{eventReader: eventReader}
}

func (h *EventHistoryHandler) Handle(ctx context.Context, qry EventHistory) (*EventHistoryResult, error) {
	pv, err := vo.NewPaymentIDVO(qry.PaymentID)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	stored, err := h.eventReader.FindByPaymentID(ctx, pv.ID)
	if err != nil {
		return nil, err
	}

	items := make([]EventHistoryItem, len(stored))
	for i, e := range stored {
		items[i] = EventHistoryItem{ID: e.ID, EventType: e.EventType, Payload: e.Payload, OccurredAt: e.OccurredAt}
	}
	return &EventHistoryResult{Events: items}, nil
}
