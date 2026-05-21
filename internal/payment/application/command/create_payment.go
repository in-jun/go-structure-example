package command

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/payment/domain"
	"github.com/in-jun/go-structure-example/internal/payment/domain/entity"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

type CreatePayment struct {
	AuctionID string
	WinnerID  string
	Amount    int64
}

type CreatePaymentResult struct {
	ID        string
	AuctionID string
	WinnerID  string
	Amount    int64
	Status    string
}

type CreatePaymentHandler struct {
	paymentRepo    domain.PaymentRepository
	eventPublisher domain.EventPublisher
	transactor     transaction.Transactor
}

func NewCreatePaymentHandler(paymentRepo domain.PaymentRepository, eventPublisher domain.EventPublisher, transactor transaction.Transactor) *CreatePaymentHandler {
	return &CreatePaymentHandler{paymentRepo: paymentRepo, eventPublisher: eventPublisher, transactor: transactor}
}

func (h *CreatePaymentHandler) Handle(ctx context.Context, cmd CreatePayment) (*CreatePaymentResult, error) {
	payment, err := entity.NewPayment(cmd.AuctionID, cmd.WinnerID, cmd.Amount)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	var result *CreatePaymentResult
	err = h.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := h.paymentRepo.Save(txCtx, payment); err != nil {
			return err
		}
		if err := h.eventPublisher.Publish(txCtx, payment.Events()...); err != nil {
			return err
		}
		payment.ClearEvents()

		result = &CreatePaymentResult{
			ID: payment.ID(), AuctionID: payment.AuctionID(),
			WinnerID: payment.WinnerID(), Amount: payment.Amount(),
			Status: payment.Status(),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
