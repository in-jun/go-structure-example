package application

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/auth/application/command"
)

type CommandUseCase interface {
	Register(ctx context.Context, cmd command.Register) error
	Login(ctx context.Context, cmd command.Login) (*command.LoginResult, error)
	Refresh(ctx context.Context, cmd command.Refresh) (*command.RefreshResult, error)
	Logout(ctx context.Context, cmd command.Logout) error
}

var _ CommandUseCase = (*service)(nil)

type service struct {
	register *command.RegisterHandler
	login    *command.LoginHandler
	refresh  *command.RefreshHandler
	logout   *command.LogoutHandler
}

func NewService(
	register *command.RegisterHandler,
	login *command.LoginHandler,
	refresh *command.RefreshHandler,
	logout *command.LogoutHandler,
) *service {
	return &service{
		register: register,
		login:    login,
		refresh:  refresh,
		logout:   logout,
	}
}

func (s *service) Register(ctx context.Context, cmd command.Register) error {
	return s.register.Handle(ctx, cmd)
}

func (s *service) Login(ctx context.Context, cmd command.Login) (*command.LoginResult, error) {
	return s.login.Handle(ctx, cmd)
}

func (s *service) Refresh(ctx context.Context, cmd command.Refresh) (*command.RefreshResult, error) {
	return s.refresh.Handle(ctx, cmd)
}

func (s *service) Logout(ctx context.Context, cmd command.Logout) error {
	return s.logout.Handle(ctx, cmd)
}
