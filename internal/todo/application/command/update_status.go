package command

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/domain"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

type UpdateStatus struct {
	UserID uint
	TodoID uint
	Status entity.Status
}

type UpdateStatusHandler struct {
	todoRepo domain.TodoRepository
}

func NewUpdateStatusHandler(todoRepo domain.TodoRepository) *UpdateStatusHandler {
	return &UpdateStatusHandler{todoRepo: todoRepo}
}

func (h *UpdateStatusHandler) Handle(ctx context.Context, cmd UpdateStatus) error {
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

	t.SetStatus(cmd.Status)
	return h.todoRepo.Update(ctx, t)
}
