package query

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

func makeResult() *entity.Todo {
	t, _ := entity.NewTodo(1, "Test", "", time.Now().Add(time.Hour))
	t.SetID(1)
	return t
}

func TestListHandler_Success(t *testing.T) {
	todos := []*entity.Todo{makeResult(), makeResult()}
	h := NewListHandler(&mockTodoRepo{todos: todos, total: 2})

	result, err := h.Handle(context.Background(), List{UserID: 1, Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if result.Total != 2 {
		t.Errorf("Total = %d, want 2", result.Total)
	}
	if len(result.Todos) != 2 {
		t.Errorf("len(Todos) = %d, want 2", len(result.Todos))
	}
}

func TestListHandler_Empty(t *testing.T) {
	h := NewListHandler(&mockTodoRepo{todos: nil, total: 0})

	result, err := h.Handle(context.Background(), List{UserID: 1, Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if result.Total != 0 {
		t.Errorf("Total = %d, want 0", result.Total)
	}
	if len(result.Todos) != 0 {
		t.Errorf("len(Todos) = %d, want 0", len(result.Todos))
	}
}
