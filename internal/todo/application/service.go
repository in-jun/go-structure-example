package application

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/todo/application/command"
	"github.com/in-jun/go-structure-example/internal/todo/application/query"
)

type CommandUseCase interface {
	Create(ctx context.Context, cmd command.Create) (*command.CreateResult, error)
	Update(ctx context.Context, cmd command.Update) error
	UpdateStatus(ctx context.Context, cmd command.UpdateStatus) error
	Delete(ctx context.Context, cmd command.Delete) error
}

type QueryUseCase interface {
	Get(ctx context.Context, qry query.Get) (*query.Result, error)
	GetList(ctx context.Context, qry query.List) (*query.ListResult, error)
}

var (
	_ CommandUseCase = (*service)(nil)
	_ QueryUseCase   = (*service)(nil)
)

type service struct {
	create       *command.CreateHandler
	update       *command.UpdateHandler
	updateStatus *command.UpdateStatusHandler
	delete       *command.DeleteHandler
	get          *query.GetHandler
	list         *query.ListHandler
}

func NewService(
	create *command.CreateHandler,
	update *command.UpdateHandler,
	updateStatus *command.UpdateStatusHandler,
	delete *command.DeleteHandler,
	get *query.GetHandler,
	list *query.ListHandler,
) *service {
	return &service{
		create:       create,
		update:       update,
		updateStatus: updateStatus,
		delete:       delete,
		get:          get,
		list:         list,
	}
}

func (s *service) Create(ctx context.Context, cmd command.Create) (*command.CreateResult, error) {
	return s.create.Handle(ctx, cmd)
}

func (s *service) Update(ctx context.Context, cmd command.Update) error {
	return s.update.Handle(ctx, cmd)
}

func (s *service) UpdateStatus(ctx context.Context, cmd command.UpdateStatus) error {
	return s.updateStatus.Handle(ctx, cmd)
}

func (s *service) Delete(ctx context.Context, cmd command.Delete) error {
	return s.delete.Handle(ctx, cmd)
}

func (s *service) Get(ctx context.Context, qry query.Get) (*query.Result, error) {
	return s.get.Handle(ctx, qry)
}

func (s *service) GetList(ctx context.Context, qry query.List) (*query.ListResult, error) {
	return s.list.Handle(ctx, qry)
}
