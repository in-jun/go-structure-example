package command

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/auth/domain"
	"github.com/in-jun/go-structure-example/internal/auth/domain/vo"
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
	v, err := vo.NewTokenStringVO(cmd.RefreshToken)
	if err != nil {
		return errors.BadRequest(err.Error())
	}

	if err := h.tokenRepo.DeleteByToken(ctx, v.Token); err != nil {
		return err
	}

	if cmd.AccessTokenJTI != "" {
		ttl := time.Duration(h.tokenGen.AccessExpirySeconds()) * time.Second
		if err2 := h.tokenRepo.BlacklistAccessToken(ctx, cmd.AccessTokenJTI, ttl); err2 != nil {
			return err2
		}
	}

	return nil
}
