package command

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/todo/domain"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
	"github.com/in-jun/go-structure-example/internal/todo/domain/vo"
)

type Create struct {
	UserID      uint
	Title       string
	Description string
	DueDate     time.Time
}

type CreateResult struct {
	ID uint
}

type CreateHandler struct {
	todoRepo domain.TodoRepository
}

func NewCreateHandler(todoRepo domain.TodoRepository) *CreateHandler {
	return &CreateHandler{todoRepo: todoRepo}
}

func (h *CreateHandler) Handle(ctx context.Context, cmd Create) (*CreateResult, error) {
	v, err := vo.NewCreateTodoVO(cmd.Title, cmd.Description, cmd.DueDate)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	t, err := entity.NewTodo(cmd.UserID, v.Title, v.Description, v.DueDate)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	if err := h.todoRepo.Save(ctx, t); err != nil {
		return nil, err
	}

	return &CreateResult{ID: t.ID()}, nil
}
