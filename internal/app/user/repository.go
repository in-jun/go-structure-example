package user

import "context"

type Repository interface {
	Save(ctx context.Context, user *User) error
	FindById(ctx context.Context, id uint) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
}
