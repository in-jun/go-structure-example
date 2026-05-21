package command

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/payment/domain"
	"github.com/in-jun/go-structure-example/internal/payment/domain/service"
	"github.com/in-jun/go-structure-example/internal/shared/query"
	"github.com/in-jun/go-structure-example/internal/payment/domain/vo"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

type RefundPayment struct {
	UserID    string
	PaymentID string
	Reason    string
}

type RefundPaymentHandler struct {
	paymentRepo    domain.PaymentRepository
	processor      *service.PaymentProcessor
	eventPublisher domain.EventPublisher
	transactor     transaction.Transactor
}

func NewRefundPaymentHandler(
	paymentRepo domain.PaymentRepository,
	processor *service.PaymentProcessor,
	eventPublisher domain.EventPublisher,
	transactor transaction.Transactor,
) *RefundPaymentHandler {
	return &RefundPaymentHandler{
		paymentRepo: paymentRepo, processor: processor,
		eventPublisher: eventPublisher, transactor: transactor,
	}
}

func (h *RefundPaymentHandler) Handle(ctx context.Context, cmd RefundPayment) error {
	pv, err := vo.NewPaymentIDVO(cmd.PaymentID)
	if err != nil {
		return errors.BadRequest(err.Error())
	}

	return h.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
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

		reason := cmd.Reason
		if reason == "" {
			reason = "User requested refund"
		}

		if err := h.processor.ProcessRefund(txCtx, payment, reason); err != nil {
			return errors.Conflict(err.Error())
		}

		if err := h.paymentRepo.Update(txCtx, payment); err != nil {
			return err
		}

		if err := h.eventPublisher.Publish(txCtx, payment.Events()...); err != nil {
			return err
		}
		payment.ClearEvents()
		return nil
	}, transaction.WithIsolation(transaction.Pessimistic))
}
