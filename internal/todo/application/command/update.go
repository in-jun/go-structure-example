package command

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/domain"
	"github.com/in-jun/go-structure-example/internal/todo/domain/vo"
)

type Update struct {
	UserID      string
	TodoID      string
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
	v, err := vo.NewUpdateTodoVO(cmd.Title, cmd.Description, cmd.DueDate)
	if err != nil {
		return errors.BadRequest(err.Error())
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

	t.Update(v.Title, v.Description, v.DueDate)
	return h.todoRepo.Update(ctx, t)
}
