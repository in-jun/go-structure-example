package http

import (
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/todo/application/query"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

func TestToTodoResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	result := &query.TodoResult{
		ID:          1,
		Title:       "Buy groceries",
		Description: "Milk and eggs",
		Status:      entity.StatusPending,
		DueDate:     now.Add(time.Hour),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	resp := toTodoResponse(result)

	if resp.ID != 1 {
		t.Errorf("ID = %d, want 1", resp.ID)
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
	result := &query.TodoListResult{
		Todos: []query.TodoResult{
			{ID: 1, Title: "Todo 1", DueDate: now.Add(time.Hour)},
			{ID: 2, Title: "Todo 2", DueDate: now.Add(2 * time.Hour)},
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
