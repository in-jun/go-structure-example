package service

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/payment/domain"
	"github.com/in-jun/go-structure-example/internal/payment/domain/entity"
)

type PaymentProcessor struct {
	gateway domain.PaymentGateway
}

func NewPaymentProcessor(gateway domain.PaymentGateway) *PaymentProcessor {
	return &PaymentProcessor{gateway: gateway}
}

func (p *PaymentProcessor) Process(ctx context.Context, payment *entity.Payment) error {
	err := p.gateway.Charge(ctx, payment.ID(), payment.Amount())
	if err != nil {
		return payment.Fail(err.Error())
	}
	return payment.Complete()
}

func (p *PaymentProcessor) ProcessRefund(ctx context.Context, payment *entity.Payment, reason string) error {
	if err := p.gateway.Refund(ctx, payment.ID(), payment.Amount()); err != nil {
		return err
	}
	return payment.Refund(reason)
}
