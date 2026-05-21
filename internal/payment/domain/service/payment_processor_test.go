package service

import (
	"context"
	"errors"
	"testing"

	"github.com/in-jun/go-structure-example/internal/payment/domain/entity"
)

type mockGatewaySuccess struct{}

func (m *mockGatewaySuccess) Charge(_ context.Context, _ string, _ int64) error { return nil }
func (m *mockGatewaySuccess) Refund(_ context.Context, _ string, _ int64) error { return nil }

type mockGatewayFail struct{}

func (m *mockGatewayFail) Charge(_ context.Context, _ string, _ int64) error {
	return errors.New("declined")
}
func (m *mockGatewayFail) Refund(_ context.Context, _ string, _ int64) error {
	return errors.New("refund failed")
}

func TestPaymentProcessor_Process_Success(t *testing.T) {
	processor := NewPaymentProcessor(&mockGatewaySuccess{})
	payment, _ := entity.NewPayment("auction-id", "winner-id", 5000)

	if err := processor.Process(context.Background(), payment); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment.Status() != entity.StatusCompleted {
		t.Errorf("expected status '%s', got '%s'", entity.StatusCompleted, payment.Status())
	}
}

func TestPaymentProcessor_Process_Failure(t *testing.T) {
	processor := NewPaymentProcessor(&mockGatewayFail{})
	payment, _ := entity.NewPayment("auction-id", "winner-id", 5000)

	if err := processor.Process(context.Background(), payment); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment.Status() != entity.StatusFailed {
		t.Errorf("expected status '%s', got '%s'", entity.StatusFailed, payment.Status())
	}
}

func TestPaymentProcessor_ProcessRefund_Success(t *testing.T) {
	processor := NewPaymentProcessor(&mockGatewaySuccess{})
	payment, _ := entity.NewPayment("auction-id", "winner-id", 5000)
	_ = processor.Process(context.Background(), payment)

	if err := processor.ProcessRefund(context.Background(), payment, "test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment.Status() != entity.StatusRefunded {
		t.Errorf("expected status '%s', got '%s'", entity.StatusRefunded, payment.Status())
	}
}

func TestPaymentProcessor_ProcessRefund_GatewayFail(t *testing.T) {
	processor := NewPaymentProcessor(&mockGatewayFail{})
	payment, _ := entity.NewPayment("auction-id", "winner-id", 5000)
	if err := payment.Complete(); err != nil {
		t.Fatal(err)
	}

	if err := processor.ProcessRefund(context.Background(), payment, "test"); err == nil {
		t.Error("expected error from gateway refund failure")
	}
}
