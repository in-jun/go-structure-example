package user

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/pkg/utils/errors"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(ctx context.Context, id uint) (*UserResponse, error) {
	user, err := s.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NotFound("User not found")
	}

	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (s *Service) Update(ctx context.Context, id uint, req UpdateUserRequest) error {
	user, err := s.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.NotFound("User not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Internal("Failed to hash password")
	}

	user.Name = req.Name
	user.Password = string(hashedPassword)

	return s.repo.Update(ctx, user)
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
