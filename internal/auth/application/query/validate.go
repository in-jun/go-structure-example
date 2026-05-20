package query

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/auth/domain"
	"github.com/in-jun/go-structure-example/internal/auth/domain/vo"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

type Validate struct {
	TokenString string
}

type Result struct {
	UserID   uint
	JTI      string
	IssuedAt int64
}

type ValidateHandler struct {
	tokenRepo domain.TokenRepository
	tokenGen  domain.TokenGenerator
}

func NewValidateHandler(tokenRepo domain.TokenRepository, tokenGen domain.TokenGenerator) *ValidateHandler {
	return &ValidateHandler{tokenRepo: tokenRepo, tokenGen: tokenGen}
}

func (h *ValidateHandler) Handle(ctx context.Context, qry Validate) (*Result, error) {
	v, err := vo.NewTokenStringVO(qry.TokenString)
	if err != nil {
		return nil, errors.Unauthorized(err.Error())
	}

	claims, err := h.tokenGen.ValidateToken(v.Token)
	if err != nil {
		return nil, errors.Unauthorized("Invalid token")
	}

	blacklisted, err := h.tokenRepo.IsAccessTokenBlacklisted(ctx, claims.JTI)
	if err != nil {
		return nil, err
	}
	if blacklisted {
		return nil, errors.Unauthorized("Token has been revoked")
	}

	revoked, err := h.tokenRepo.IsRevokedForUser(ctx, claims.UserID, claims.IssuedAt)
	if err != nil {
		return nil, err
	}
	if revoked {
		return nil, errors.Unauthorized("Token has been revoked")
	}

	return &Result{UserID: claims.UserID, JTI: claims.JTI, IssuedAt: claims.IssuedAt}, nil
}
