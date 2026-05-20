package command

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

func TestUpdateHandler_Success(t *testing.T) {
	todo := makeTodo()
	h := NewUpdateHandler(&mockTodoRepo{todo: todo})

	err := h.Handle(context.Background(), Update{
		UserID:  testUUID,
		TodoID:  todo.ID(),
		Title:   "Updated Title",
		DueDate: time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestUpdateHandler_NotFound(t *testing.T) {
	h := NewUpdateHandler(&mockTodoRepo{err: errors.NotFound("not found")})

	err := h.Handle(context.Background(), Update{
		UserID:  testUUID,
		TodoID:  "nonexistent-id",
		Title:   "Updated",
		DueDate: time.Now().Add(time.Hour),
	})
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestUpdateHandler_Forbidden(t *testing.T) {
	todo := makeTodo()
	h := NewUpdateHandler(&mockTodoRepo{todo: todo})

	err := h.Handle(context.Background(), Update{
		UserID:  "660e8400-e29b-41d4-a716-446655440000",
		TodoID:  todo.ID(),
		Title:   "Updated",
		DueDate: time.Now().Add(time.Hour),
	})
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}
}

func TestUpdateHandler_EmptyTitle(t *testing.T) {
	todo := makeTodo()
	h := NewUpdateHandler(&mockTodoRepo{todo: todo})

	err := h.Handle(context.Background(), Update{
		UserID:  testUUID,
		TodoID:  todo.ID(),
		Title:   "",
		DueDate: time.Now().Add(time.Hour),
	})
	if err == nil {
		t.Fatal("expected error for empty title, got nil")
	}
}
