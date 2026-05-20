package command

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/user/domain"
	"github.com/in-jun/go-structure-example/internal/user/domain/vo"
)

type UpdatePassword struct {
	UserID          uint
	CurrentPassword string
	NewPassword     string
}

type UpdatePasswordHandler struct {
	userRepo       domain.UserRepository
	passwordHasher domain.PasswordHasher
}

func NewUpdatePasswordHandler(userRepo domain.UserRepository, passwordHasher domain.PasswordHasher) *UpdatePasswordHandler {
	return &UpdatePasswordHandler{userRepo: userRepo, passwordHasher: passwordHasher}
}

func (h *UpdatePasswordHandler) Handle(ctx context.Context, cmd UpdatePassword) error {
	v, err := vo.NewUpdatePasswordVO(cmd.CurrentPassword, cmd.NewPassword)
	if err != nil {
		return errors.BadRequest(err.Error())
	}

	u, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	if u == nil {
		return errors.NotFound("User not found")
	}

	if !h.passwordHasher.Compare(u.HashedPassword(), v.CurrentPassword) {
		return errors.Unauthorized("Current password is incorrect")
	}

	hashed, err := h.passwordHasher.Hash(v.NewPassword)
	if err != nil {
		return errors.Internal("Failed to hash password")
	}

	u.SetPassword(hashed)
	return h.userRepo.Update(ctx, u)
}
