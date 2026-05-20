package application

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/todo/application/command"
	"github.com/in-jun/go-structure-example/internal/todo/application/query"
)

type CommandUseCase interface {
	Create(ctx context.Context, cmd command.Create) (*command.CreateResult, error)
	Update(ctx context.Context, cmd command.Update) error
	UpdateStatus(ctx context.Context, cmd command.UpdateStatus) error
	Delete(ctx context.Context, cmd command.Delete) error
}

type QueryUseCase interface {
	GetTodo(ctx context.Context, qry query.GetTodo) (*query.TodoResult, error)
	ListTodos(ctx context.Context, qry query.ListTodos) (*query.TodoListResult, error)
}

var (
	_ CommandUseCase = (*service)(nil)
	_ QueryUseCase   = (*service)(nil)
)

type service struct {
	create       *command.CreateHandler
	update       *command.UpdateHandler
	updateStatus *command.UpdateStatusHandler
	delete       *command.DeleteHandler
	getTodo      *query.GetTodoHandler
	listTodos    *query.ListTodosHandler
}

func NewService(
	create *command.CreateHandler,
	update *command.UpdateHandler,
	updateStatus *command.UpdateStatusHandler,
	delete *command.DeleteHandler,
	getTodo *query.GetTodoHandler,
	listTodos *query.ListTodosHandler,
) *service {
	return &service{
		create:       create,
		update:       update,
		updateStatus: updateStatus,
		delete:       delete,
		getTodo:      getTodo,
		listTodos:    listTodos,
	}
}

func (s *service) Create(ctx context.Context, cmd command.Create) (*command.CreateResult, error) {
	return s.create.Handle(ctx, cmd)
}

func (s *service) Update(ctx context.Context, cmd command.Update) error {
	return s.update.Handle(ctx, cmd)
}

func (s *service) UpdateStatus(ctx context.Context, cmd command.UpdateStatus) error {
	return s.updateStatus.Handle(ctx, cmd)
}

func (s *service) Delete(ctx context.Context, cmd command.Delete) error {
	return s.delete.Handle(ctx, cmd)
}

func (s *service) GetTodo(ctx context.Context, qry query.GetTodo) (*query.TodoResult, error) {
	return s.getTodo.Handle(ctx, qry)
}

func (s *service) ListTodos(ctx context.Context, qry query.ListTodos) (*query.TodoListResult, error) {
	return s.listTodos.Handle(ctx, qry)
}
