package command

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

type mockTodoRepo struct {
	todo *entity.Todo
	err  error
}

func (m *mockTodoRepo) Save(_ context.Context, t *entity.Todo) error {
	t.SetID(1)
	return m.err
}
func (m *mockTodoRepo) FindByID(_ context.Context, _ uint) (*entity.Todo, error) {
	return m.todo, m.err
}
func (m *mockTodoRepo) FindByUserID(_ context.Context, _ uint, _, _ int) ([]*entity.Todo, int64, error) {
	return nil, 0, m.err
}
func (m *mockTodoRepo) Update(_ context.Context, _ *entity.Todo) error { return m.err }
func (m *mockTodoRepo) Delete(_ context.Context, _ uint) error         { return m.err }

func TestCreateHandler_Success(t *testing.T) {
	h := NewCreateHandler(&mockTodoRepo{})

	result, err := h.Handle(context.Background(), Create{
		UserID:  1,
		Title:   "Buy groceries",
		DueDate: time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if result.ID == 0 {
		t.Error("expected non-zero ID after save")
	}
}

func TestCreateHandler_EmptyTitle(t *testing.T) {
	h := NewCreateHandler(&mockTodoRepo{})

	_, err := h.Handle(context.Background(), Create{
		UserID:  1,
		Title:   "",
		DueDate: time.Now().Add(time.Hour),
	})
	if err == nil {
		t.Fatal("expected error for empty title, got nil")
	}
}

func TestCreateHandler_ZeroUserID(t *testing.T) {
	h := NewCreateHandler(&mockTodoRepo{})

	_, err := h.Handle(context.Background(), Create{
		UserID:  0,
		Title:   "Valid Title",
		DueDate: time.Now().Add(time.Hour),
	})
	if err == nil {
		t.Fatal("expected error for zero user ID, got nil")
	}
}
