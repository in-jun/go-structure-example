package domain

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/user/domain/entity"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) bool
}
