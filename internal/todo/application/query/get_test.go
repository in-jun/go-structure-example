package query

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

type mockTodoRepo struct {
	todo  *entity.Todo
	todos []*entity.Todo
	total int64
	err   error
}

func (m *mockTodoRepo) Save(_ context.Context, _ *entity.Todo) error { return m.err }
func (m *mockTodoRepo) FindByID(_ context.Context, _ uint) (*entity.Todo, error) {
	return m.todo, m.err
}
func (m *mockTodoRepo) FindByUserID(_ context.Context, _ uint, _, _ int) ([]*entity.Todo, int64, error) {
	return m.todos, m.total, m.err
}
func (m *mockTodoRepo) Update(_ context.Context, _ *entity.Todo) error { return m.err }
func (m *mockTodoRepo) Delete(_ context.Context, _ uint) error         { return m.err }

func makeTodo() *entity.Todo {
	t, _ := entity.NewTodo(1, "Test Todo", "description", time.Now().Add(time.Hour))
	t.SetID(1)
	return t
}

func TestGetHandler_Success(t *testing.T) {
	h := NewGetHandler(&mockTodoRepo{todo: makeTodo()})

	result, err := h.Handle(context.Background(), Get{UserID: 1, TodoID: 1})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if result.Title != "Test Todo" {
		t.Errorf("Title = %q, want %q", result.Title, "Test Todo")
	}
}

func TestGetHandler_NotFound(t *testing.T) {
	h := NewGetHandler(&mockTodoRepo{err: errors.NotFound("not found")})

	_, err := h.Handle(context.Background(), Get{UserID: 1, TodoID: 99})
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestGetHandler_Forbidden(t *testing.T) {
	h := NewGetHandler(&mockTodoRepo{todo: makeTodo()})

	_, err := h.Handle(context.Background(), Get{UserID: 999, TodoID: 1})
	if err == nil {
		t.Fatal("expected forbidden error for wrong user, got nil")
	}
}
