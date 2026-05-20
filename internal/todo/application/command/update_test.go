package command

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

func TestUpdateHandler_Success(t *testing.T) {
	h := NewUpdateHandler(&mockTodoRepo{todo: makeTodo()})

	err := h.Handle(context.Background(), Update{
		UserID:  1,
		TodoID:  1,
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
		UserID:  1,
		TodoID:  99,
		Title:   "Updated",
		DueDate: time.Now().Add(time.Hour),
	})
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestUpdateHandler_Forbidden(t *testing.T) {
	h := NewUpdateHandler(&mockTodoRepo{todo: makeTodo()})

	err := h.Handle(context.Background(), Update{
		UserID:  999,
		TodoID:  1,
		Title:   "Updated",
		DueDate: time.Now().Add(time.Hour),
	})
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}
}

func TestUpdateHandler_EmptyTitle(t *testing.T) {
	h := NewUpdateHandler(&mockTodoRepo{todo: makeTodo()})

	err := h.Handle(context.Background(), Update{
		UserID:  1,
		TodoID:  1,
		Title:   "",
		DueDate: time.Now().Add(time.Hour),
	})
	if err == nil {
		t.Fatal("expected error for empty title, got nil")
	}
}
