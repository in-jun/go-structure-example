package command

import (
	"context"
	"log/slog"
	"time"

	"github.com/in-jun/go-structure-example/internal/auth/domain"
	"github.com/in-jun/go-structure-example/internal/auth/domain/entity"
	"github.com/in-jun/go-structure-example/internal/auth/domain/vo"
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
	v, err := vo.NewRefreshTokenVO(cmd.RefreshToken)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	old, err := h.tokenRepo.FindByToken(ctx, v.Token)
	if err != nil {
		return nil, err
	}
	if old == nil {
		userID, err := h.tokenRepo.FindUsedToken(ctx, v.Token)
		if err != nil {
			return nil, err
		}
		if userID != "" {
			if err := h.tokenRepo.DeleteByUserID(ctx, userID); err != nil {
				slog.Warn("failed to revoke all user tokens", "error", err)
			}
			return nil, errors.Unauthorized("Token reuse detected, all sessions revoked")
		}
		return nil, errors.Unauthorized("Invalid refresh token")
	}

	if old.IsExpired() {
		if err := h.tokenRepo.DeleteByToken(ctx, v.Token); err != nil {
			slog.Warn("failed to delete expired refresh token", "error", err)
		}
		return nil, errors.Unauthorized("Refresh token expired")
	}

	if err := h.tokenRepo.DeleteByToken(ctx, v.Token); err != nil {
		return nil, err
	}
	if err := h.tokenRepo.MarkTokenUsed(ctx, v.Token, old.UserID()); err != nil {
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
