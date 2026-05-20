package command

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/user/domain"
)

type Delete struct {
	UserID uint
}

type DeleteHandler struct {
	userRepo domain.UserRepository
}

func NewDeleteHandler(userRepo domain.UserRepository) *DeleteHandler {
	return &DeleteHandler{userRepo: userRepo}
}

func (h *DeleteHandler) Handle(ctx context.Context, cmd Delete) error {
	return h.userRepo.Delete(ctx, cmd.UserID)
}
