package query

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/domain"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

type GetTodo struct {
	UserID uint
	TodoID uint
}

type TodoResult struct {
	ID          uint
	Title       string
	Description string
	Status      entity.Status
	DueDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GetTodoHandler struct {
	todoRepo domain.TodoRepository
}

func NewGetTodoHandler(todoRepo domain.TodoRepository) *GetTodoHandler {
	return &GetTodoHandler{todoRepo: todoRepo}
}

func (h *GetTodoHandler) Handle(ctx context.Context, qry GetTodo) (*TodoResult, error) {
	t, err := h.todoRepo.FindByID(ctx, qry.TodoID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errors.NotFound("Todo not found")
	}
	if t.UserID() != qry.UserID {
		return nil, errors.Forbidden("Not authorized to access this todo")
	}

	return toTodoResult(t), nil
}

func toTodoResult(t *entity.Todo) *TodoResult {
	return &TodoResult{
		ID:          t.ID(),
		Title:       t.Title(),
		Description: t.Description(),
		Status:      t.Status(),
		DueDate:     t.DueDate(),
		CreatedAt:   t.CreatedAt(),
		UpdatedAt:   t.UpdatedAt(),
	}
}
