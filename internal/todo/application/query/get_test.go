package query

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/domain"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

const testUUID = "550e8400-e29b-41d4-a716-446655440000"

type mockTodoRepo struct {
	todo  *entity.Todo
	todos []*entity.Todo
	total int64
	err   error
}

func (m *mockTodoRepo) Save(_ context.Context, _ *entity.Todo) error { return m.err }
func (m *mockTodoRepo) FindByID(_ context.Context, _ string) (*entity.Todo, error) {
	return m.todo, m.err
}
func (m *mockTodoRepo) FindByUserID(_ context.Context, _ string, _, _ int) ([]*entity.Todo, int64, error) {
	return m.todos, m.total, m.err
}
func (m *mockTodoRepo) Update(_ context.Context, _ *entity.Todo) error { return m.err }
func (m *mockTodoRepo) Delete(_ context.Context, _ string) error       { return m.err }

var _ domain.TodoRepository = (*mockTodoRepo)(nil)

func makeTodo() *entity.Todo {
	t, _ := entity.NewTodo(testUUID, "Test Todo", "description", time.Now().Add(time.Hour))
	return t
}

func TestGetHandler_Success(t *testing.T) {
	todo := makeTodo()
	h := NewGetHandler(&mockTodoRepo{todo: todo})

	result, err := h.Handle(context.Background(), Get{UserID: testUUID, TodoID: todo.ID()})
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
	if result.Title != "Test Todo" {
		t.Errorf("Title = %q, want %q", result.Title, "Test Todo")
	}
}

func TestGetHandler_NotFound(t *testing.T) {
	h := NewGetHandler(&mockTodoRepo{err: errors.NotFound("not found")})

	_, err := h.Handle(context.Background(), Get{UserID: testUUID, TodoID: "some-id"})
	if err == nil {
		t.Fatal("expected error for not found, got nil")
	}
}

func TestGetHandler_Forbidden(t *testing.T) {
	todo := makeTodo()
	h := NewGetHandler(&mockTodoRepo{todo: todo})

	_, err := h.Handle(context.Background(), Get{UserID: "660e8400-e29b-41d4-a716-446655440000", TodoID: todo.ID()})
	if err == nil {
		t.Fatal("expected forbidden error for wrong user, got nil")
	}
}
