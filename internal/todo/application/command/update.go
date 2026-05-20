package command

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/domain"
)

type Update struct {
	UserID      uint
	TodoID      uint
	Title       string
	Description string
	DueDate     time.Time
}

type UpdateHandler struct {
	todoRepo domain.TodoRepository
}

func NewUpdateHandler(todoRepo domain.TodoRepository) *UpdateHandler {
	return &UpdateHandler{todoRepo: todoRepo}
}

func (h *UpdateHandler) Handle(ctx context.Context, cmd Update) error {
	if cmd.Title == "" {
		return errors.BadRequest("Title is required")
	}

	t, err := h.todoRepo.FindByID(ctx, cmd.TodoID)
	if err != nil {
		return err
	}
	if t == nil {
		return errors.NotFound("Todo not found")
	}
	if t.UserID() != cmd.UserID {
		return errors.Forbidden("Not authorized to update this todo")
	}

	t.Update(cmd.Title, cmd.Description, cmd.DueDate)
	return h.todoRepo.Update(ctx, t)
}
