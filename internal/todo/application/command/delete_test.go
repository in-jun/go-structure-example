package command

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

func makeTodo() *entity.Todo {
	t, _ := entity.NewTodo(testUUID, "Test", "", time.Now().Add(time.Hour))
	return t
}

func TestDeleteHandler_Success(t *testing.T) {
	todo := makeTodo()
	h := NewDeleteHandler(&mockTodoRepo{todo: todo})

	err := h.Handle(context.Background(), Delete{UserID: testUUID, TodoID: todo.ID()})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestDeleteHandler_NotFound(t *testing.T) {
	h := NewDeleteHandler(&mockTodoRepo{err: errors.NotFound("not found")})

	err := h.Handle(context.Background(), Delete{UserID: testUUID, TodoID: "some-id"})
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestDeleteHandler_Forbidden(t *testing.T) {
	todo := makeTodo()
	h := NewDeleteHandler(&mockTodoRepo{todo: todo})

	err := h.Handle(context.Background(), Delete{UserID: "660e8400-e29b-41d4-a716-446655440000", TodoID: todo.ID()})
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}
}
