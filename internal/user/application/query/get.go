package query

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/user/domain"
)

type GetUser struct {
	UserID uint
}

type UserResult struct {
	ID        uint
	Email     string
	Name      string
	CreatedAt time.Time
}

type GetUserHandler struct {
	userRepo domain.UserRepository
}

func NewGetUserHandler(userRepo domain.UserRepository) *GetUserHandler {
	return &GetUserHandler{userRepo: userRepo}
}

func (h *GetUserHandler) Handle(ctx context.Context, qry GetUser) (*UserResult, error) {
	u, err := h.userRepo.FindByID(ctx, qry.UserID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.NotFound("User not found")
	}

	return &UserResult{
		ID:        u.ID(),
		Email:     u.Email(),
		Name:      u.Name(),
		CreatedAt: u.CreatedAt(),
	}, nil
}
