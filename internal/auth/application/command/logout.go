package command

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/auth/domain"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

type Logout struct {
	RefreshToken   string
	AccessTokenJTI string
}

type LogoutHandler struct {
	tokenRepo domain.TokenRepository
	tokenGen  domain.TokenGenerator
}

func NewLogoutHandler(tokenRepo domain.TokenRepository, tokenGen domain.TokenGenerator) *LogoutHandler {
	return &LogoutHandler{tokenRepo: tokenRepo, tokenGen: tokenGen}
}

func (h *LogoutHandler) Handle(ctx context.Context, cmd Logout) error {
	if cmd.RefreshToken == "" {
		return errors.BadRequest("Refresh token is required")
	}

	if err := h.tokenRepo.DeleteByToken(ctx, cmd.RefreshToken); err != nil {
		return err
	}

	if cmd.AccessTokenJTI != "" {
		ttl := time.Duration(h.tokenGen.AccessExpirySeconds()) * time.Second
		if err := h.tokenRepo.BlacklistAccessToken(ctx, cmd.AccessTokenJTI, ttl); err != nil {
			return err
		}
	}

	return nil
}
