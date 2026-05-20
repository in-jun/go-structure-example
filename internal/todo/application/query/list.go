package query

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/todo/domain"
)

type List struct {
	UserID string
	Page   int
	Limit  int
}

type ListResult struct {
	Todos []Result
	Total int64
}

type ListHandler struct {
	todoRepo domain.TodoRepository
}

func NewListHandler(todoRepo domain.TodoRepository) *ListHandler {
	return &ListHandler{todoRepo: todoRepo}
}

func (h *ListHandler) Handle(ctx context.Context, qry List) (*ListResult, error) {
	todos, total, err := h.todoRepo.FindByUserID(ctx, qry.UserID, qry.Page, qry.Limit)
	if err != nil {
		return nil, err
	}

	results := make([]Result, 0, len(todos))
	for _, t := range todos {
		results = append(results, *toTodoResult(t))
	}

	return &ListResult{Todos: results, Total: total}, nil
}
