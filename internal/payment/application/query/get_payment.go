package query

import (
	"context"
	"time"

	"github.com/in-jun/go-structure-example/internal/payment/domain"
	"github.com/in-jun/go-structure-example/internal/payment/domain/vo"
	"github.com/in-jun/go-structure-example/internal/shared/errors"
)

type GetPayment struct {
	PaymentID string
}

type Result struct {
	ID        string
	AuctionID string
	WinnerID  string
	Amount    int64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GetPaymentHandler struct {
	paymentRepo domain.PaymentRepository
}

func NewGetPaymentHandler(paymentRepo domain.PaymentRepository) *GetPaymentHandler {
	return &GetPaymentHandler{paymentRepo: paymentRepo}
}

func (h *GetPaymentHandler) Handle(ctx context.Context, qry GetPayment) (*Result, error) {
	pv, err := vo.NewPaymentIDVO(qry.PaymentID)
	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	payment, err := h.paymentRepo.FindByID(ctx, pv.ID)
	if err != nil {
		return nil, err
	}
	if payment == nil {
		return nil, errors.NotFound("Payment not found")
	}

	return &Result{
		ID: payment.ID(), AuctionID: payment.AuctionID(),
		WinnerID: payment.WinnerID(), Amount: payment.Amount(),
		Status: payment.Status(), CreatedAt: payment.CreatedAt(),
		UpdatedAt: payment.UpdatedAt(),
	}, nil
}
