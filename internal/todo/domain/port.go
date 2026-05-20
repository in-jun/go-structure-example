package domain

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

type TodoRepository interface {
	Save(ctx context.Context, todo *entity.Todo) error
	FindByID(ctx context.Context, id uint) (*entity.Todo, error)
	FindByUserID(ctx context.Context, userID uint, page, limit int) ([]*entity.Todo, int64, error)
	Update(ctx context.Context, todo *entity.Todo) error
	Delete(ctx context.Context, id uint) error
}
