package command

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/payment/domain"
	"github.com/in-jun/go-structure-example/internal/payment/domain/entity"
	"github.com/in-jun/go-structure-example/internal/shared/query"
	"github.com/in-jun/go-structure-example/internal/payment/domain/service"
	"github.com/in-jun/go-structure-example/internal/payment/domain/vo"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

type ConfirmPayment struct {
	UserID    string
	PaymentID string
}

type ConfirmPaymentHandler struct {
	paymentRepo    domain.PaymentRepository
	processor      *service.PaymentProcessor
	eventPublisher domain.EventPublisher
	transactor     transaction.Transactor
}

func NewConfirmPaymentHandler(
	paymentRepo domain.PaymentRepository,
	processor *service.PaymentProcessor,
	eventPublisher domain.EventPublisher,
	transactor transaction.Transactor,
) *ConfirmPaymentHandler {
	return &ConfirmPaymentHandler{
		paymentRepo: paymentRepo, processor: processor,
		eventPublisher: eventPublisher, transactor: transactor,
	}
}

func (h *ConfirmPaymentHandler) Handle(ctx context.Context, cmd ConfirmPayment) error {
	pv, err := vo.NewPaymentIDVO(cmd.PaymentID)
	if err != nil {
		return errors.BadRequest(err.Error())
	}

	var declined bool
	err = h.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		payment, err := h.paymentRepo.FindByID(txCtx, pv.ID, query.ForUpdate())
		if err != nil {
			return err
		}
		if payment == nil {
			return errors.NotFound("Payment not found")
		}
		if !payment.IsOwnedBy(cmd.UserID) {
			return errors.Forbidden("Not authorized")
		}

		if err := h.processor.Process(txCtx, payment); err != nil {
			return errors.Conflict(err.Error())
		}

		if err := h.paymentRepo.Update(txCtx, payment); err != nil {
			return err
		}

		if err := h.eventPublisher.Publish(txCtx, payment.Events()...); err != nil {
			return err
		}
		payment.ClearEvents()

		declined = payment.Status() == entity.StatusFailed
		return nil
	}, transaction.WithIsolation(transaction.Pessimistic))
	if err != nil {
		return err
	}
	if declined {
		return errors.BadRequest("Payment was declined")
	}
	return nil
}
