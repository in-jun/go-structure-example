package auth

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/app/user"
	"github.com/in-jun/go-structure-example/internal/pkg/config"
	"github.com/in-jun/go-structure-example/internal/pkg/utils/errors"
	"github.com/in-jun/go-structure-example/internal/pkg/utils/jwt"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) error
	Login(ctx context.Context, req LoginRequest) (*AuthResponse, error)
	Refresh(ctx context.Context, req RefreshRequest) (*AuthResponse, error)
	Logout(ctx context.Context, userID uint, refreshToken, jti string) error
}

var _ Service = (*service)(nil)

type service struct {
	authRepo Repository
	userRepo user.Repository
}

func NewService(authRepo Repository, userRepo user.Repository) Service {
	return &service{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) error {
	existing, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.Conflict("Email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Internal("Failed to hash password")
	}

	return s.userRepo.Save(ctx, &user.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	})
}

func (s *service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	u, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.Unauthorized("Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, errors.Unauthorized("Invalid credentials")
	}

	return s.generateTokens(ctx, u.ID)
}

func (s *service) Refresh(ctx context.Context, req RefreshRequest) (*AuthResponse, error) {
	old, err := s.authRepo.FindByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if old == nil {
		return nil, errors.Unauthorized("Invalid refresh token")
	}

	if time.Now().After(old.ExpiresAt()) {
		return nil, errors.Unauthorized("Refresh token expired")
	}

	if err := s.authRepo.DeleteByToken(ctx, req.RefreshToken); err != nil {
		return nil, err
	}

	return s.generateTokens(ctx, old.UserID())
}

func (s *service) Logout(ctx context.Context, userID uint, refreshToken, jti string) error {
	if err := s.authRepo.DeleteByToken(ctx, refreshToken); err != nil {
		return err
	}

	expiry, err := time.ParseDuration(config.AppConfig.JWTAccessExpiry)
	if err != nil {
		return errors.Internal("Invalid access token expiry configuration")
	}

	return s.authRepo.BlacklistAccessToken(ctx, jti, expiry)
}

func (s *service) generateTokens(ctx context.Context, userID uint) (*AuthResponse, error) {
	accessToken, err := jwt.GenerateAccessToken(userID)
	if err != nil {
		return nil, errors.Internal("Failed to generate access token")
	}

	refreshExpiry, err := time.ParseDuration(config.AppConfig.JWTRefreshExpiry)
	if err != nil {
		return nil, errors.Internal("Invalid refresh token expiry configuration")
	}

	rt, err := NewRefreshToken(userID, time.Now().Add(refreshExpiry))
	if err != nil {
		return nil, errors.Internal("Failed to generate refresh token")
	}

	if err := s.authRepo.Save(ctx, rt); err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: rt.Token(),
		ExpiresIn:    jwt.AccessExpirySeconds(),
	}, nil
}
