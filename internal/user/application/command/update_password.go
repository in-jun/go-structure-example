package command

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/user/domain"
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
	if cmd.CurrentPassword == "" || cmd.NewPassword == "" {
		return errors.BadRequest("Current and new password are required")
	}
	if len(cmd.NewPassword) < 6 {
		return errors.BadRequest("New password must be at least 6 characters")
	}

	u, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	if u == nil {
		return errors.NotFound("User not found")
	}

	if !h.passwordHasher.Compare(u.HashedPassword(), cmd.CurrentPassword) {
		return errors.Unauthorized("Current password is incorrect")
	}

	hashed, err := h.passwordHasher.Hash(cmd.NewPassword)
	if err != nil {
		return errors.Internal("Failed to hash password")
	}

	u.SetPassword(hashed)
	return h.userRepo.Update(ctx, u)
}
