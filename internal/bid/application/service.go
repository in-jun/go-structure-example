package application

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/bid/application/command"
	"github.com/in-jun/go-structure-example/internal/bid/application/query"
)

type CommandUseCase interface {
	PlaceBid(ctx context.Context, cmd command.PlaceBid) (*command.PlaceBidResult, error)
	DetermineWinner(ctx context.Context, cmd command.DetermineWinner) error
}

type QueryUseCase interface {
	GetHighest(ctx context.Context, qry query.GetHighest) (*query.Result, error)
	ListBids(ctx context.Context, qry query.ListBids) (*query.ListResult, error)
}

var (
	_ CommandUseCase = (*service)(nil)
	_ QueryUseCase   = (*service)(nil)
)

type service struct {
	placeBid        *command.PlaceBidHandler
	determineWinner *command.DetermineWinnerHandler
	getHighest      *query.GetHighestHandler
	listBids        *query.ListBidsHandler
}

func NewService(
	placeBid *command.PlaceBidHandler,
	determineWinner *command.DetermineWinnerHandler,
	getHighest *query.GetHighestHandler,
	listBids *query.ListBidsHandler,
) *service {
	return &service{
		placeBid: placeBid, determineWinner: determineWinner,
		getHighest: getHighest, listBids: listBids,
	}
}

func (s *service) PlaceBid(ctx context.Context, cmd command.PlaceBid) (*command.PlaceBidResult, error) {
	return s.placeBid.Handle(ctx, cmd)
}
func (s *service) DetermineWinner(ctx context.Context, cmd command.DetermineWinner) error {
	return s.determineWinner.Handle(ctx, cmd)
}
func (s *service) GetHighest(ctx context.Context, qry query.GetHighest) (*query.Result, error) {
	return s.getHighest.Handle(ctx, qry)
}
func (s *service) ListBids(ctx context.Context, qry query.ListBids) (*query.ListResult, error) {
	return s.listBids.Handle(ctx, qry)
}
