package query

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/domain"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

type Get struct {
	UserID uint
	TodoID uint
}

type Result struct {
	ID          uint
	Title       string
	Description string
	Status      entity.Status
	DueDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GetHandler struct {
	todoRepo domain.TodoRepository
}

func NewGetHandler(todoRepo domain.TodoRepository) *GetHandler {
	return &GetHandler{todoRepo: todoRepo}
}

func (h *GetHandler) Handle(ctx context.Context, qry Get) (*Result, error) {
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

func toTodoResult(t *entity.Todo) *Result {
	return &Result{
		ID:          t.ID(),
		Title:       t.Title(),
		Description: t.Description(),
		Status:      t.Status(),
		DueDate:     t.DueDate(),
		CreatedAt:   t.CreatedAt(),
		UpdatedAt:   t.UpdatedAt(),
	}
}
