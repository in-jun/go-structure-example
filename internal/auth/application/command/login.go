package command

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/auth/domain"
	"github.com/in-jun/go-structure-example/internal/auth/domain/entity"
	"github.com/in-jun/go-structure-example/internal/auth/domain/vo"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

type Login struct {
	Email    string
	Password string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
}

type LoginHandler struct {
	userRepo       domain.UserRepository
	tokenRepo      domain.TokenRepository
	tokenGen       domain.TokenGenerator
	passwordHasher domain.PasswordHasher
}

func NewLoginHandler(
	userRepo domain.UserRepository,
	tokenRepo domain.TokenRepository,
	tokenGen domain.TokenGenerator,
	passwordHasher domain.PasswordHasher,
) *LoginHandler {
	return &LoginHandler{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		tokenGen:       tokenGen,
		passwordHasher: passwordHasher,
	}
}

func (h *LoginHandler) Handle(ctx context.Context, cmd Login) (*LoginResult, error) {
	v, err := vo.NewLoginVO(cmd.Email, cmd.Password)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	user, err := h.userRepo.FindByEmail(ctx, v.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.Unauthorized("Invalid credentials")
	}

	if !h.passwordHasher.Compare(user.HashedPassword(), v.Password) {
		return nil, errors.Unauthorized("Invalid credentials")
	}

	accessToken, err := h.tokenGen.GenerateAccessToken(user.ID())
	if err != nil {
		return nil, errors.Internal("Failed to generate access token")
	}

	rt, err := entity.NewRefreshToken(user.ID(), time.Now().Add(h.tokenGen.RefreshExpiry()))
	if err != nil {
		return nil, errors.Internal("Failed to create refresh token")
	}

	if err := h.tokenRepo.Save(ctx, rt); err != nil {
		return nil, err
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: rt.Token(),
		ExpiresIn:    h.tokenGen.AccessExpirySeconds(),
	}, nil
}
