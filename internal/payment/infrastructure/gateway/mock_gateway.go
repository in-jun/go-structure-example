package gateway

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/in-jun/go-structure-example/internal/payment/domain"
)

var _ domain.PaymentGateway = (*MockGateway)(nil)

type MockGateway struct{}

func NewMockGateway() domain.PaymentGateway {
	return &MockGateway{}
}

func rollPercent() int64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(100))
	return n.Int64()
}

func (g *MockGateway) Charge(_ context.Context, _ string, _ int64) error {
	if rollPercent() < 10 {
		return errors.New("payment gateway: transaction declined")
	}
	return nil
}

func (g *MockGateway) Refund(_ context.Context, _ string, _ int64) error {
	if rollPercent() < 5 {
		return errors.New("payment gateway: refund failed")
	}
	return nil
}
