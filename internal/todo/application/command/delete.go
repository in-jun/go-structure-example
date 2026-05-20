package command

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/domain"
)

type Delete struct {
	UserID string
	TodoID string
}

type DeleteHandler struct {
	todoRepo domain.TodoRepository
}

func NewDeleteHandler(todoRepo domain.TodoRepository) *DeleteHandler {
	return &DeleteHandler{todoRepo: todoRepo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd Delete) error {
	t, err := h.todoRepo.FindByID(ctx, cmd.TodoID)
	if err != nil {
		return err
	}
	if t == nil {
		return errors.NotFound("Todo not found")
	}
	if t.UserID() != cmd.UserID {
		return errors.Forbidden("Not authorized to delete this todo")
	}

	return h.todoRepo.Delete(ctx, cmd.TodoID)
}
