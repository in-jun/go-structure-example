package application

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/payment/application/command"
	"github.com/in-jun/go-structure-example/internal/payment/application/query"
)

type CommandUseCase interface {
	CreatePayment(ctx context.Context, cmd command.CreatePayment) (*command.CreatePaymentResult, error)
	ConfirmPayment(ctx context.Context, cmd command.ConfirmPayment) error
	RefundPayment(ctx context.Context, cmd command.RefundPayment) error
}

type QueryUseCase interface {
	GetPayment(ctx context.Context, qry query.GetPayment) (*query.Result, error)
	GetEvents(ctx context.Context, qry query.EventHistory) (*query.EventHistoryResult, error)
}

var (
	_ CommandUseCase = (*service)(nil)
	_ QueryUseCase   = (*service)(nil)
)

type service struct {
	createPayment  *command.CreatePaymentHandler
	confirmPayment *command.ConfirmPaymentHandler
	refundPayment  *command.RefundPaymentHandler
	getPayment     *query.GetPaymentHandler
	getEvents      *query.EventHistoryHandler
}

func NewService(
	createPayment *command.CreatePaymentHandler,
	confirmPayment *command.ConfirmPaymentHandler,
	refundPayment *command.RefundPaymentHandler,
	getPayment *query.GetPaymentHandler,
	getEvents *query.EventHistoryHandler,
) *service {
	return &service{
		createPayment: createPayment, confirmPayment: confirmPayment,
		refundPayment: refundPayment, getPayment: getPayment, getEvents: getEvents,
	}
}

func (s *service) CreatePayment(ctx context.Context, cmd command.CreatePayment) (*command.CreatePaymentResult, error) {
	return s.createPayment.Handle(ctx, cmd)
}
func (s *service) ConfirmPayment(ctx context.Context, cmd command.ConfirmPayment) error {
	return s.confirmPayment.Handle(ctx, cmd)
}
func (s *service) RefundPayment(ctx context.Context, cmd command.RefundPayment) error {
	return s.refundPayment.Handle(ctx, cmd)
}
func (s *service) GetPayment(ctx context.Context, qry query.GetPayment) (*query.Result, error) {
	return s.getPayment.Handle(ctx, qry)
}
func (s *service) GetEvents(ctx context.Context, qry query.EventHistory) (*query.EventHistoryResult, error) {
	return s.getEvents.Handle(ctx, qry)
}
