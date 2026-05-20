package application

import (
	"context"
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/application/command"
	"github.com/in-jun/go-structure-example/internal/todo/application/query"
	"github.com/in-jun/go-structure-example/internal/todo/domain"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

type mockTodoRepo struct {
	todo  *entity.Todo
	todos []*entity.Todo
	total int64
	err   error
}

func (m *mockTodoRepo) Save(_ context.Context, t *entity.Todo) error {
	t.SetID(1)
	return m.err
}
func (m *mockTodoRepo) FindByID(_ context.Context, _ uint) (*entity.Todo, error) {
	return m.todo, m.err
}
func (m *mockTodoRepo) FindByUserID(_ context.Context, _ uint, _, _ int) ([]*entity.Todo, int64, error) {
	return m.todos, m.total, m.err
}
func (m *mockTodoRepo) Update(_ context.Context, _ *entity.Todo) error { return m.err }
func (m *mockTodoRepo) Delete(_ context.Context, _ uint) error         { return m.err }

func newTestService(repo *mockTodoRepo) *service {
	return NewService(
		command.NewCreateHandler(repo),
		command.NewUpdateHandler(repo),
		command.NewUpdateStatusHandler(repo),
		command.NewDeleteHandler(repo),
		query.NewGetHandler(repo),
		query.NewListHandler(repo),
	)
}

func makeTodo() *entity.Todo {
	t, _ := entity.NewTodo(1, "Test Todo", "description", time.Now().Add(time.Hour))
	t.SetID(1)
	return t
}

func TestTodoService_Create(t *testing.T) {
	svc := newTestService(&mockTodoRepo{})

	result, err := svc.Create(context.Background(), command.Create{
		UserID:  1,
		Title:   "Buy groceries",
		DueDate: time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if result.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestTodoService_Get(t *testing.T) {
	svc := newTestService(&mockTodoRepo{todo: makeTodo()})

	result, err := svc.Get(context.Background(), query.Get{UserID: 1, TodoID: 1})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if result.Title != "Test Todo" {
		t.Errorf("expected Test Todo, got %q", result.Title)
	}
}

func TestTodoService_Get_NotFound(t *testing.T) {
	svc := newTestService(&mockTodoRepo{err: errors.NotFound("not found")})

	_, err := svc.Get(context.Background(), query.Get{UserID: 1, TodoID: 99})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTodoService_Get_Forbidden(t *testing.T) {
	todo := makeTodo()
	svc := newTestService(&mockTodoRepo{todo: todo})

	_, err := svc.Get(context.Background(), query.Get{UserID: 999, TodoID: 1})
	if err == nil {
		t.Fatal("expected forbidden error, got nil")
	}
}

func TestTodoService_GetList(t *testing.T) {
	svc := newTestService(&mockTodoRepo{todos: []*entity.Todo{makeTodo()}, total: 1})

	result, err := svc.GetList(context.Background(), query.List{UserID: 1, Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("GetList() error = %v", err)
	}
	if result.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Total)
	}
}

func TestTodoService_Update(t *testing.T) {
	svc := newTestService(&mockTodoRepo{todo: makeTodo()})

	err := svc.Update(context.Background(), command.Update{
		UserID:  1,
		TodoID:  1,
		Title:   "Updated",
		DueDate: time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
}

func TestTodoService_UpdateStatus(t *testing.T) {
	svc := newTestService(&mockTodoRepo{todo: makeTodo()})

	err := svc.UpdateStatus(context.Background(), command.UpdateStatus{
		UserID: 1,
		TodoID: 1,
		Status: entity.StatusCompleted,
	})
	if err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}
}

func TestTodoService_Delete(t *testing.T) {
	svc := newTestService(&mockTodoRepo{todo: makeTodo()})

	err := svc.Delete(context.Background(), command.Delete{UserID: 1, TodoID: 1})
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
}

var _ domain.TodoRepository = (*mockTodoRepo)(nil)
