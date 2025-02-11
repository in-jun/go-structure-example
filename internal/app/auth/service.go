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

type Service struct {
	authRepo Repository
	userRepo user.Repository
}

func NewService(authRepo Repository, userRepo user.Repository) *Service {
	return &Service{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) error {
	existingUser, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.BadRequest("Email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Internal("Failed to hash password")
	}

	newUser := &user.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	return s.userRepo.Save(ctx, newUser)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.Unauthorized("Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.Unauthorized("Invalid credentials")
	}

	return s.generateTokens(ctx, user.ID)
}

func (s *Service) Refresh(ctx context.Context, req RefreshRequest) (*AuthResponse, error) {
	oldToken, err := s.authRepo.FindByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	if oldToken == nil {
		return nil, errors.Unauthorized("Invalid refresh token")
	}

	if time.Now().After(oldToken.ExpiresAt) {
		s.authRepo.DeleteByToken(ctx, req.RefreshToken)
		return nil, errors.Unauthorized("Refresh token expired")
	}

	if err := s.authRepo.DeleteByToken(ctx, req.RefreshToken); err != nil {
		return nil, err
	}

	return s.generateTokens(ctx, oldToken.UserID)
}

func (s *Service) Logout(ctx context.Context, userID uint) error {
	return s.authRepo.DeleteByUserId(ctx, userID)
}

func (s *Service) generateTokens(ctx context.Context, userID uint) (*AuthResponse, error) {
	accessToken, err := jwt.GenerateAccessToken(userID)
	if err != nil {
		return nil, errors.Internal("Failed to generate access token")
	}

	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return nil, errors.Internal("Failed to generate refresh token")
	}

	duration, err := time.ParseDuration(config.AppConfig.JWTRefreshExpiry)
	if err != nil {
		return nil, errors.Internal("Invalid refresh token expiry duration")
	}

	if err := s.authRepo.Save(ctx, &RefreshToken{
		RefreshToken: refreshToken,
		UserID:       userID,
		ExpiresAt:    time.Now().Add(duration),
	}); err != nil {
		return nil, err
	}

	accessExpiry, _ := time.ParseDuration(config.AppConfig.JWTAccessExpiry)
	expiresIn := int(accessExpiry.Seconds())

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}
