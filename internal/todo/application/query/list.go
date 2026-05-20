package query

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/todo/domain"
)

type ListTodos struct {
	UserID uint
	Page   int
	Limit  int
}

type TodoListResult struct {
	Todos []TodoResult
	Total int64
}

type ListTodosHandler struct {
	todoRepo domain.TodoRepository
}

func NewListTodosHandler(todoRepo domain.TodoRepository) *ListTodosHandler {
	return &ListTodosHandler{todoRepo: todoRepo}
}

func (h *ListTodosHandler) Handle(ctx context.Context, qry ListTodos) (*TodoListResult, error) {
	todos, total, err := h.todoRepo.FindByUserID(ctx, qry.UserID, qry.Page, qry.Limit)
	if err != nil {
		return nil, err
	}

	results := make([]TodoResult, 0, len(todos))
	for _, t := range todos {
		results = append(results, *toTodoResult(t))
	}

	return &TodoListResult{Todos: results, Total: total}, nil
}
