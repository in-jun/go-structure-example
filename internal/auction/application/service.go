package application

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/auction/application/command"
	"github.com/in-jun/go-structure-example/internal/auction/application/query"
)

type CommandUseCase interface {
	Create(ctx context.Context, cmd command.Create) (*command.CreateResult, error)
	Open(ctx context.Context, cmd command.Open) error
	Close(ctx context.Context, cmd command.Close) error
	Settle(ctx context.Context, cmd command.Settle) error
	Cancel(ctx context.Context, cmd command.Cancel) error
}

type QueryUseCase interface {
	GetByID(ctx context.Context, qry query.Get) (*query.Result, error)
	GetList(ctx context.Context, qry query.List) (*query.ListResult, error)
	GetEvents(ctx context.Context, qry query.EventHistory) (*query.EventHistoryResult, error)
}

var (
	_ CommandUseCase = (*service)(nil)
	_ QueryUseCase   = (*service)(nil)
)

type service struct {
	create       *command.CreateHandler
	open         *command.OpenHandler
	close        *command.CloseHandler
	settle       *command.SettleHandler
	cancel       *command.CancelHandler
	get          *query.GetHandler
	list         *query.ListHandler
	eventHistory *query.EventHistoryHandler
}

func NewService(
	create *command.CreateHandler,
	open *command.OpenHandler,
	close *command.CloseHandler,
	settle *command.SettleHandler,
	cancel *command.CancelHandler,
	get *query.GetHandler,
	list *query.ListHandler,
	eventHistory *query.EventHistoryHandler,
) *service {
	return &service{
		create: create, open: open, close: close,
		settle: settle, cancel: cancel,
		get: get, list: list, eventHistory: eventHistory,
	}
}

func (s *service) Create(ctx context.Context, cmd command.Create) (*command.CreateResult, error) {
	return s.create.Handle(ctx, cmd)
}
func (s *service) Open(ctx context.Context, cmd command.Open) error {
	return s.open.Handle(ctx, cmd)
}
func (s *service) Close(ctx context.Context, cmd command.Close) error {
	return s.close.Handle(ctx, cmd)
}
func (s *service) Settle(ctx context.Context, cmd command.Settle) error {
	return s.settle.Handle(ctx, cmd)
}
func (s *service) Cancel(ctx context.Context, cmd command.Cancel) error {
	return s.cancel.Handle(ctx, cmd)
}
func (s *service) GetByID(ctx context.Context, qry query.Get) (*query.Result, error) {
	return s.get.Handle(ctx, qry)
}
func (s *service) GetList(ctx context.Context, qry query.List) (*query.ListResult, error) {
	return s.list.Handle(ctx, qry)
}
func (s *service) GetEvents(ctx context.Context, qry query.EventHistory) (*query.EventHistoryResult, error) {
	return s.eventHistory.Handle(ctx, qry)
}
