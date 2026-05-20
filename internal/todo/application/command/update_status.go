package command

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/domain"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
	"github.com/in-jun/go-structure-example/internal/todo/domain/vo"
)

type UpdateStatus struct {
	UserID string
	TodoID string
	Status entity.Status
}

type UpdateStatusHandler struct {
	todoRepo domain.TodoRepository
}

func NewUpdateStatusHandler(todoRepo domain.TodoRepository) *UpdateStatusHandler {
	return &UpdateStatusHandler{todoRepo: todoRepo}
}

func (h *UpdateStatusHandler) Handle(ctx context.Context, cmd UpdateStatus) error {
	v, err := vo.NewUpdateStatusVO(cmd.Status)
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

	t.SetStatus(v.Status)
	return h.todoRepo.Update(ctx, t)
}
