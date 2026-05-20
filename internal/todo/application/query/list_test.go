package query

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

func makeResult() *entity.Todo {
	t, _ := entity.NewTodo(testUUID, "Test", "", time.Now().Add(time.Hour))
	return t
}

func TestListHandler_Success(t *testing.T) {
	todos := []*entity.Todo{makeResult(), makeResult()}
	h := NewListHandler(&mockTodoRepo{todos: todos, total: 2})

	result, err := h.Handle(context.Background(), List{UserID: testUUID, Page: 1, Limit: 10})
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

	result, err := h.Handle(context.Background(), List{UserID: testUUID, Page: 1, Limit: 10})
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

func TestListHandler_RepositoryError(t *testing.T) {
	h := NewListHandler(&mockTodoRepo{err: errors.Internal("db error")})

	_, err := h.Handle(context.Background(), List{UserID: testUUID, Page: 1, Limit: 10})
	if err == nil {
		t.Fatal("expected error for repository failure, got nil")
	}
}
