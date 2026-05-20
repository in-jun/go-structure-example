package user

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/pkg/utils/errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetByID(ctx context.Context, id uint) (*UserResponse, error)
	UpdateProfile(ctx context.Context, id uint, req UpdateProfileRequest) error
	UpdatePassword(ctx context.Context, id uint, req UpdatePasswordRequest) error
	Delete(ctx context.Context, id uint) error
}

var _ Service = (*service)(nil)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByID(ctx context.Context, id uint) (*UserResponse, error) {
	u, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.NotFound("User not found")
	}

	return &UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
	}, nil
}

func (s *service) UpdateProfile(ctx context.Context, id uint, req UpdateProfileRequest) error {
	u, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if u == nil {
		return errors.NotFound("User not found")
	}

	u.Name = req.Name
	return s.repo.Update(ctx, u)
}

func (s *service) UpdatePassword(ctx context.Context, id uint, req UpdatePasswordRequest) error {
	u, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if u == nil {
		return errors.NotFound("User not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.Unauthorized("Current password is incorrect")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.Internal("Failed to hash password")
	}

	u.Password = string(hashed)
	return s.repo.Update(ctx, u)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
