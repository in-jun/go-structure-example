package command

import (
	"context"
	"testing"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

func TestUpdateStatusHandler_Success(t *testing.T) {
	todo := makeTodo()
	h := NewUpdateStatusHandler(&mockTodoRepo{todo: todo})

	err := h.Handle(context.Background(), UpdateStatus{
		UserID: testUUID,
		TodoID: todo.ID(),
		Status: entity.StatusCompleted,
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestUpdateStatusHandler_InvalidStatus(t *testing.T) {
	todo := makeTodo()
	h := NewUpdateStatusHandler(&mockTodoRepo{todo: todo})

	err := h.Handle(context.Background(), UpdateStatus{
		UserID: testUUID,
		TodoID: todo.ID(),
		Status: entity.Status("invalid"),
	})
	if err == nil {
		t.Fatal("expected error for invalid status, got nil")
	}
}

func TestUpdateStatusHandler_NotFound(t *testing.T) {
	h := NewUpdateStatusHandler(&mockTodoRepo{err: errors.NotFound("not found")})

	err := h.Handle(context.Background(), UpdateStatus{
		UserID: testUUID,
		TodoID: "nonexistent-id",
		Status: entity.StatusCompleted,
	})
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestUpdateStatusHandler_Forbidden(t *testing.T) {
	todo := makeTodo()
	h := NewUpdateStatusHandler(&mockTodoRepo{todo: todo})

	err := h.Handle(context.Background(), UpdateStatus{
		UserID: "660e8400-e29b-41d4-a716-446655440000",
		TodoID: todo.ID(),
		Status: entity.StatusCompleted,
	})
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}
}
