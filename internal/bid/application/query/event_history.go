package query

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/bid/domain"
	"github.com/in-jun/go-structure-example/internal/bid/domain/vo"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

type EventHistory struct {
	AuctionID string
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
	av, err := vo.NewAuctionIDVO(qry.AuctionID)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	stored, err := h.eventReader.FindByAuctionID(ctx, av.ID)
	if err != nil {
		return nil, err
	}

	items := make([]EventHistoryItem, len(stored))
	for i, e := range stored {
		items[i] = EventHistoryItem{ID: e.ID, EventType: e.EventType, Payload: e.Payload, OccurredAt: e.OccurredAt}
	}
	return &EventHistoryResult{Events: items}, nil
}
