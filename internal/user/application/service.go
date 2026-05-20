package application

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/user/application/command"
	"github.com/in-jun/go-structure-example/internal/user/application/query"
)

type CommandUseCase interface {
	UpdateProfile(ctx context.Context, cmd command.UpdateProfile) error
	UpdatePassword(ctx context.Context, cmd command.UpdatePassword) error
	Delete(ctx context.Context, cmd command.Delete) error
}

type QueryUseCase interface {
	GetUser(ctx context.Context, qry query.GetUser) (*query.UserResult, error)
}

var (
	_ CommandUseCase = (*service)(nil)
	_ QueryUseCase   = (*service)(nil)
)

type service struct {
	updateProfile   *command.UpdateProfileHandler
	updatePassword  *command.UpdatePasswordHandler
	delete          *command.DeleteHandler
	getUser         *query.GetUserHandler
}

func NewService(
	updateProfile *command.UpdateProfileHandler,
	updatePassword *command.UpdatePasswordHandler,
	delete *command.DeleteHandler,
	getUser *query.GetUserHandler,
) *service {
	return &service{
		updateProfile:  updateProfile,
		updatePassword: updatePassword,
		delete:         delete,
		getUser:        getUser,
	}
}

func (s *service) UpdateProfile(ctx context.Context, cmd command.UpdateProfile) error {
	return s.updateProfile.Handle(ctx, cmd)
}

func (s *service) UpdatePassword(ctx context.Context, cmd command.UpdatePassword) error {
	return s.updatePassword.Handle(ctx, cmd)
}

func (s *service) Delete(ctx context.Context, cmd command.Delete) error {
	return s.delete.Handle(ctx, cmd)
}

func (s *service) GetUser(ctx context.Context, qry query.GetUser) (*query.UserResult, error) {
	return s.getUser.Handle(ctx, qry)
}
