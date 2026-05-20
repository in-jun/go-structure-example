package command

import (
	"context"

	apperrors "github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/user/domain"
)

type UpdateProfile struct {
	UserID uint
	Name   string
}

type UpdateProfileHandler struct {
	userRepo domain.UserRepository
}

func NewUpdateProfileHandler(userRepo domain.UserRepository) *UpdateProfileHandler {
	return &UpdateProfileHandler{userRepo: userRepo}
}

func (h *UpdateProfileHandler) Handle(ctx context.Context, cmd UpdateProfile) error {
	if cmd.Name == "" {
		return apperrors.BadRequest("Name is required")
	}

	u, err := h.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	if u == nil {
		return apperrors.NotFound("User not found")
	}

	u.SetName(cmd.Name)
	return h.userRepo.Update(ctx, u)
}
