package application

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/in-jun/go-structure-example/internal/payment/application/command"
	"github.com/in-jun/go-structure-example/internal/payment/application/query"
	"github.com/in-jun/go-structure-example/internal/payment/domain/entity"
	sharedQuery "github.com/in-jun/go-structure-example/internal/shared/query"
	domainEvent "github.com/in-jun/go-structure-example/internal/payment/domain/event"
	domainService "github.com/in-jun/go-structure-example/internal/payment/domain/service"
	"github.com/in-jun/go-structure-example/internal/shared/transaction"
)

type mockPaymentRepo struct {
	payment *entity.Payment
	err     error
}

func (m *mockPaymentRepo) Save(_ context.Context, _ *entity.Payment) error { return m.err }
func (m *mockPaymentRepo) FindByID(_ context.Context, _ string, _ ...sharedQuery.Option) (*entity.Payment, error) {
	return m.payment, m.err
}
func (m *mockPaymentRepo) Update(_ context.Context, _ *entity.Payment) error { return m.err }

type mockGateway struct{}

func (m *mockGateway) Charge(_ context.Context, _ string, _ int64) error { return nil }
func (m *mockGateway) Refund(_ context.Context, _ string, _ int64) error { return nil }

type mockPublisher struct{}

func (m *mockPublisher) Publish(_ context.Context, _ ...domainEvent.Event) error { return nil }

type mockTransactor struct{}

func (m *mockTransactor) WithinTransaction(_ context.Context, fn func(ctx context.Context) error, _ ...transaction.TxOption) error {
	return fn(context.Background())
}

func newTestService(repo *mockPaymentRepo) *service {
	processor := domainService.NewPaymentProcessor(&mockGateway{})
	return NewService(
		command.NewCreatePaymentHandler(repo, &mockPublisher{}, &mockTransactor{}),
		command.NewConfirmPaymentHandler(repo, processor, &mockPublisher{}, &mockTransactor{}),
		command.NewRefundPaymentHandler(repo, processor, &mockPublisher{}, &mockTransactor{}),
		query.NewGetPaymentHandler(repo),
	)
}

func TestPaymentService_CreatePayment(t *testing.T) {
	svc := newTestService(&mockPaymentRepo{})

	result, err := svc.CreatePayment(context.Background(), command.CreatePayment{
		AuctionID: uuid.New().String(),
		WinnerID:  uuid.New().String(),
		Amount:    5000,
	})
	if err != nil {
		t.Fatalf("CreatePayment() error = %v", err)
	}
	if result.Status != entity.StatusPending {
		t.Errorf("Status = %q, want %q", result.Status, entity.StatusPending)
	}
}

func TestPaymentService_ConfirmPayment(t *testing.T) {
	winnerID := uuid.New().String()
	payment, _ := entity.NewPayment(uuid.New().String(), winnerID, 5000)
	repo := &mockPaymentRepo{payment: payment}
	svc := newTestService(repo)

	err := svc.ConfirmPayment(context.Background(), command.ConfirmPayment{
		UserID:    winnerID,
		PaymentID: payment.ID(),
	})
	if err != nil {
		t.Fatalf("ConfirmPayment() error = %v", err)
	}
}

func TestPaymentService_ConfirmPayment_NotOwner(t *testing.T) {
	winnerID := uuid.New().String()
	payment, _ := entity.NewPayment(uuid.New().String(), winnerID, 5000)
	repo := &mockPaymentRepo{payment: payment}
	svc := newTestService(repo)

	err := svc.ConfirmPayment(context.Background(), command.ConfirmPayment{
		UserID:    uuid.New().String(),
		PaymentID: payment.ID(),
	})
	if err == nil {
		t.Error("expected error for non-owner")
	}
}

func TestPaymentService_GetPayment(t *testing.T) {
	now := time.Now()
	payment := entity.ReconstructPayment(uuid.New().String(), uuid.New().String(), uuid.New().String(), 5000, entity.StatusPending, now, now)
	svc := newTestService(&mockPaymentRepo{payment: payment})

	result, err := svc.GetPayment(context.Background(), query.GetPayment{PaymentID: payment.ID()})
	if err != nil {
		t.Fatalf("GetPayment() error = %v", err)
	}
	if result.Amount != 5000 {
		t.Errorf("Amount = %d, want 5000", result.Amount)
	}
}
