package domain

import (
	"context"

	"github.com/in-jun/go-structure-example/internal/payment/domain/entity"
	"github.com/in-jun/go-structure-example/internal/payment/domain/event"
	"github.com/in-jun/go-structure-example/internal/shared/query"
)

type PaymentRepository interface {
	Save(ctx context.Context, payment *entity.Payment) error
	FindByID(ctx context.Context, id string, opts ...query.Option) (*entity.Payment, error)
	Update(ctx context.Context, payment *entity.Payment) error
}

type PaymentGateway interface {
	Charge(ctx context.Context, paymentID string, amount int64) error
	Refund(ctx context.Context, paymentID string, amount int64) error
}

type EventPublisher interface {
	Publish(ctx context.Context, events ...event.Event) error
}

type EventReader interface {
	FindByPaymentID(ctx context.Context, paymentID string) ([]event.StoredEvent, error)
}
