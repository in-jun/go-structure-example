package command

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/auth/domain"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

type LogoutAll struct {
	UserID uint
}

type LogoutAllHandler struct {
	tokenRepo domain.TokenRepository
	tokenGen  domain.TokenGenerator
}

func NewLogoutAllHandler(tokenRepo domain.TokenRepository, tokenGen domain.TokenGenerator) *LogoutAllHandler {
	return &LogoutAllHandler{tokenRepo: tokenRepo, tokenGen: tokenGen}
}

func (h *LogoutAllHandler) Handle(ctx context.Context, cmd LogoutAll) error {
	if cmd.UserID == 0 {
		return errors.BadRequest("User ID is required")
	}

	if err := h.tokenRepo.DeleteByUserID(ctx, cmd.UserID); err != nil {
		return err
	}

	ttl := time.Duration(h.tokenGen.AccessExpirySeconds()) * time.Second
	return h.tokenRepo.RevokeAllAccessTokens(ctx, cmd.UserID, ttl)
}
