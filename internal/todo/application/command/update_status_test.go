package command

import (
	"context"
	"testing"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

func TestUpdateStatusHandler_Success(t *testing.T) {
	h := NewUpdateStatusHandler(&mockTodoRepo{todo: makeTodo()})

	err := h.Handle(context.Background(), UpdateStatus{
		UserID: 1,
		TodoID: 1,
		Status: entity.StatusCompleted,
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestUpdateStatusHandler_InvalidStatus(t *testing.T) {
	h := NewUpdateStatusHandler(&mockTodoRepo{todo: makeTodo()})

	err := h.Handle(context.Background(), UpdateStatus{
		UserID: 1,
		TodoID: 1,
		Status: entity.Status("invalid"),
	})
	if err == nil {
		t.Fatal("expected error for invalid status, got nil")
	}
}

func TestUpdateStatusHandler_NotFound(t *testing.T) {
	h := NewUpdateStatusHandler(&mockTodoRepo{err: errors.NotFound("not found")})

	err := h.Handle(context.Background(), UpdateStatus{
		UserID: 1,
		TodoID: 99,
		Status: entity.StatusCompleted,
	})
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestUpdateStatusHandler_Forbidden(t *testing.T) {
	h := NewUpdateStatusHandler(&mockTodoRepo{todo: makeTodo()})

	err := h.Handle(context.Background(), UpdateStatus{
		UserID: 999,
		TodoID: 1,
		Status: entity.StatusCompleted,
	})
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}
}
