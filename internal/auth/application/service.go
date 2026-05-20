package application

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/auth/application/command"
	"github.com/in-jun/go-structure-example/internal/auth/application/query"
)

type CommandUseCase interface {
	Register(ctx context.Context, cmd command.Register) error
	Login(ctx context.Context, cmd command.Login) (*command.LoginResult, error)
	Refresh(ctx context.Context, cmd command.Refresh) (*command.RefreshResult, error)
	Logout(ctx context.Context, cmd command.Logout) error
	LogoutAll(ctx context.Context, cmd command.LogoutAll) error
}

type QueryUseCase interface {
	ValidateToken(ctx context.Context, qry query.Validate) (*query.Result, error)
}

var (
	_ CommandUseCase = (*service)(nil)
	_ QueryUseCase   = (*service)(nil)
)

type service struct {
	register  *command.RegisterHandler
	login     *command.LoginHandler
	refresh   *command.RefreshHandler
	logout    *command.LogoutHandler
	logoutAll *command.LogoutAllHandler
	validate  *query.ValidateHandler
}

func NewService(
	register *command.RegisterHandler,
	login *command.LoginHandler,
	refresh *command.RefreshHandler,
	logout *command.LogoutHandler,
	logoutAll *command.LogoutAllHandler,
	validate *query.ValidateHandler,
) *service {
	return &service{
		register:  register,
		login:     login,
		refresh:   refresh,
		logout:    logout,
		logoutAll: logoutAll,
		validate:  validate,
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

func (s *service) LogoutAll(ctx context.Context, cmd command.LogoutAll) error {
	return s.logoutAll.Handle(ctx, cmd)
}

func (s *service) ValidateToken(ctx context.Context, qry query.Validate) (*query.Result, error) {
	return s.validate.Handle(ctx, qry)
}
