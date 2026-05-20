package query

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/user/domain"
)

type Get struct {
	UserID string
}

type Result struct {
	ID        string
	Email     string
	Name      string
	CreatedAt time.Time
}

type GetHandler struct {
	userRepo domain.UserRepository
}

func NewGetHandler(userRepo domain.UserRepository) *GetHandler {
	return &GetHandler{userRepo: userRepo}
}

func (h *GetHandler) Handle(ctx context.Context, qry Get) (*Result, error) {
	u, err := h.userRepo.FindByID(ctx, qry.UserID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.NotFound("User not found")
	}

	return &Result{
		ID:        u.ID(),
		Email:     u.Email(),
		Name:      u.Name(),
		CreatedAt: u.CreatedAt(),
	}, nil
}
