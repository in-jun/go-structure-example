package command

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/auth/domain"
	"github.com/in-jun/go-structure-example/internal/auth/domain/entity"
	"github.com/in-jun/go-structure-example/internal/auth/domain/vo"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

type Register struct {
	Email    string
	Password string
	Name     string
}

type RegisterHandler struct {
	userRepo       domain.UserRepository
	passwordHasher domain.PasswordHasher
}

func NewRegisterHandler(userRepo domain.UserRepository, passwordHasher domain.PasswordHasher) *RegisterHandler {
	return &RegisterHandler{userRepo: userRepo, passwordHasher: passwordHasher}
}

func (h *RegisterHandler) Handle(ctx context.Context, cmd Register) error {
	v, err := vo.NewRegisterVO(cmd.Email, cmd.Password, cmd.Name)
	if err != nil {
		return errors.BadRequest(err.Error())
	}

	existing, err := h.userRepo.FindByEmail(ctx, v.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.Conflict("Email already registered")
	}

	hashedPassword, err := h.passwordHasher.Hash(v.Password)
	if err != nil {
		return errors.Internal("Failed to hash password")
	}

	user, err := entity.NewUser(v.Email, hashedPassword, v.Name)
	if err != nil {
		return errors.Internal("Failed to create user")
	}

	return h.userRepo.Save(ctx, user)
}
