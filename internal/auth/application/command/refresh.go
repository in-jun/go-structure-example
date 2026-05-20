package command

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/auth/domain"
	"github.com/in-jun/go-structure-example/internal/auth/domain/entity"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

type Refresh struct {
	RefreshToken string
}

type RefreshResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

type RefreshHandler struct {
	tokenRepo domain.TokenRepository
	tokenGen  domain.TokenGenerator
}

func NewRefreshHandler(tokenRepo domain.TokenRepository, tokenGen domain.TokenGenerator) *RefreshHandler {
	return &RefreshHandler{tokenRepo: tokenRepo, tokenGen: tokenGen}
}

func (h *RefreshHandler) Handle(ctx context.Context, cmd Refresh) (*RefreshResult, error) {
	if cmd.RefreshToken == "" {
		return nil, errors.BadRequest("Refresh token is required")
	}

	old, err := h.tokenRepo.FindByToken(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, err
	}
	if old == nil {
		return nil, errors.Unauthorized("Invalid refresh token")
	}

	if old.IsExpired() {
		_ = h.tokenRepo.DeleteByToken(ctx, cmd.RefreshToken)
		return nil, errors.Unauthorized("Refresh token expired")
	}

	if err := h.tokenRepo.DeleteByToken(ctx, cmd.RefreshToken); err != nil {
		return nil, err
	}

	accessToken, err := h.tokenGen.GenerateAccessToken(old.UserID())
	if err != nil {
		return nil, errors.Internal("Failed to generate access token")
	}

	rt, err := entity.NewRefreshToken(old.UserID(), time.Now().Add(h.tokenGen.RefreshExpiry()))
	if err != nil {
		return nil, errors.Internal("Failed to create refresh token")
	}

	if err := h.tokenRepo.Save(ctx, rt); err != nil {
		return nil, err
	}

	return &RefreshResult{
		AccessToken:  accessToken,
		RefreshToken: rt.Token(),
		ExpiresIn:    h.tokenGen.AccessExpirySeconds(),
	}, nil
}
