package todo

import "context"

type Repository interface {
	Save(ctx context.Context, todo *Todo) error
	Update(ctx context.Context, todo *Todo) error
	Delete(ctx context.Context, id uint) error
	FindById(ctx context.Context, id uint) (*Todo, error)
	FindByUserId(ctx context.Context, userID uint, page, limit int) ([]*Todo, int64, error)
	UpdateStatus(ctx context.Context, id uint, status Status) error
}
