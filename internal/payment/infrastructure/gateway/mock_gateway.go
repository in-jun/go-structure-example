package gateway

import (
	"context"
	"errors"
	"math/rand"

	"github.com/in-jun/go-structure-example/internal/payment/domain"
)

var _ domain.PaymentGateway = (*MockGateway)(nil)

type MockGateway struct{}

func NewMockGateway() domain.PaymentGateway {
	return &MockGateway{}
}

func (g *MockGateway) Charge(_ context.Context, _ string, _ int64) error {
	if rand.Float64() < 0.1 {
		return errors.New("payment gateway: transaction declined")
	}
	return nil
}

func (g *MockGateway) Refund(_ context.Context, _ string, _ int64) error {
	if rand.Float64() < 0.05 {
		return errors.New("payment gateway: refund failed")
	}
	return nil
}
