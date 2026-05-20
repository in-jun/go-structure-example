package command

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/todo/domain"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

const testUUID = "550e8400-e29b-41d4-a716-446655440000"

type mockTodoRepo struct {
	todo *entity.Todo
	err  error
}

func (m *mockTodoRepo) Save(_ context.Context, _ *entity.Todo) error { return m.err }
func (m *mockTodoRepo) FindByID(_ context.Context, _ string) (*entity.Todo, error) {
	return m.todo, m.err
}
func (m *mockTodoRepo) FindByUserID(_ context.Context, _ string, _, _ int) ([]*entity.Todo, int64, error) {
	return nil, 0, m.err
}
func (m *mockTodoRepo) Update(_ context.Context, _ *entity.Todo) error { return m.err }
func (m *mockTodoRepo) Delete(_ context.Context, _ string) error       { return m.err }

var _ domain.TodoRepository = (*mockTodoRepo)(nil)

func TestCreateHandler_Success(t *testing.T) {
	h := NewCreateHandler(&mockTodoRepo{})

	result, err := h.Handle(context.Background(), Create{
		UserID:  testUUID,
		Title:   "Buy groceries",
		DueDate: time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if result.ID == "" {
		t.Error("expected non-empty ID after save")
	}
}

func TestCreateHandler_EmptyTitle(t *testing.T) {
	h := NewCreateHandler(&mockTodoRepo{})

	_, err := h.Handle(context.Background(), Create{
		UserID:  testUUID,
		Title:   "",
		DueDate: time.Now().Add(time.Hour),
	})
	if err == nil {
		t.Fatal("expected error for empty title, got nil")
	}
}

func TestCreateHandler_EmptyUserID(t *testing.T) {
	h := NewCreateHandler(&mockTodoRepo{})

	_, err := h.Handle(context.Background(), Create{
		UserID:  "",
		Title:   "Valid Title",
		DueDate: time.Now().Add(time.Hour),
	})
	if err == nil {
		t.Fatal("expected error for empty user ID, got nil")
	}
}
