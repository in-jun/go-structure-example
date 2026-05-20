package http

import (
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/todo/application/query"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

const otherUUID = "660e8400-e29b-41d4-a716-446655440000"

func TestToTodoResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	result := &query.Result{
		ID:          testUUID,
		Title:       "Buy groceries",
		Description: "Milk and eggs",
		Status:      entity.StatusPending,
		DueDate:     now.Add(time.Hour),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	resp := toTodoResponse(result)

	if resp.ID != testUUID {
		t.Errorf("ID = %q, want %q", resp.ID, testUUID)
	}
	if resp.Title != "Buy groceries" {
		t.Errorf("Title = %q, want %q", resp.Title, "Buy groceries")
	}
	if resp.Status != entity.StatusPending {
		t.Errorf("Status = %q, want %q", resp.Status, entity.StatusPending)
	}
}

func TestToTodoListResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	result := &query.ListResult{
		Todos: []query.Result{
			{ID: testUUID, Title: "Todo 1", DueDate: now.Add(time.Hour)},
			{ID: otherUUID, Title: "Todo 2", DueDate: now.Add(2 * time.Hour)},
		},
		Total: 2,
	}

	resp := toTodoListResponse(result)

	if resp.Total != 2 {
		t.Errorf("Total = %d, want 2", resp.Total)
	}
	if len(resp.Todos) != 2 {
		t.Errorf("len(Todos) = %d, want 2", len(resp.Todos))
	}
	if resp.Todos[0].Title != "Todo 1" {
		t.Errorf("Todos[0].Title = %q, want %q", resp.Todos[0].Title, "Todo 1")
	}
}
