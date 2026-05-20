package command

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

func makeTodo() *entity.Todo {
	t, _ := entity.NewTodo(1, "Test", "", time.Now().Add(time.Hour))
	t.SetID(1)
	return t
}

func TestDeleteHandler_Success(t *testing.T) {
	h := NewDeleteHandler(&mockTodoRepo{todo: makeTodo()})

	err := h.Handle(context.Background(), Delete{UserID: 1, TodoID: 1})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestDeleteHandler_NotFound(t *testing.T) {
	h := NewDeleteHandler(&mockTodoRepo{err: errors.NotFound("not found")})

	err := h.Handle(context.Background(), Delete{UserID: 1, TodoID: 99})
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestDeleteHandler_Forbidden(t *testing.T) {
	h := NewDeleteHandler(&mockTodoRepo{todo: makeTodo()})

	err := h.Handle(context.Background(), Delete{UserID: 999, TodoID: 1})
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}
}
