package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

var (
	testAuctionID = uuid.New().String()
	testWinnerID  = uuid.New().String()
)

func TestNewPayment(t *testing.T) {
	payment, err := NewPayment(testAuctionID, testWinnerID, 5000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := uuid.Parse(payment.ID()); err != nil {
		t.Errorf("expected valid UUID, got '%s'", payment.ID())
	}
	if payment.Status() != StatusPending {
		t.Errorf("expected status '%s', got '%s'", StatusPending, payment.Status())
	}
	if payment.Amount() != 5000 {
		t.Errorf("expected amount 5000, got %d", payment.Amount())
	}
}

func TestNewPayment_Invariants(t *testing.T) {
	tests := []struct {
		name      string
		auctionID string
		winnerID  string
	}{
		{"empty auctionID", "", testWinnerID},
		{"empty winnerID", testAuctionID, ""},
		{"both empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPayment(tt.auctionID, tt.winnerID, 5000)
			if err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}

func TestPayment_Complete(t *testing.T) {
	payment, _ := NewPayment(testAuctionID, testWinnerID, 5000)
	payment.ClearEvents()

	if err := payment.Complete(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment.Status() != StatusCompleted {
		t.Errorf("expected status '%s', got '%s'", StatusCompleted, payment.Status())
	}
	if len(payment.Events()) != 1 {
		t.Errorf("expected 1 event, got %d", len(payment.Events()))
	}

	if err := payment.Complete(); err == nil {
		t.Error("expected error completing already completed payment")
	}
}

func TestPayment_Fail(t *testing.T) {
	payment, _ := NewPayment(testAuctionID, testWinnerID, 5000)
	payment.ClearEvents()

	if err := payment.Fail("declined"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment.Status() != StatusFailed {
		t.Errorf("expected status '%s', got '%s'", StatusFailed, payment.Status())
	}
	if len(payment.Events()) != 1 {
		t.Errorf("expected 1 event, got %d", len(payment.Events()))
	}
}

func TestPayment_Fail_NotPending(t *testing.T) {
	payment, _ := NewPayment(testAuctionID, testWinnerID, 5000)
	if err := payment.Complete(); err != nil {
		t.Fatal(err)
	}

	if err := payment.Fail("reason"); err == nil {
		t.Error("expected error failing completed payment")
	}
}

func TestPayment_IsOwnedBy(t *testing.T) {
	payment, _ := NewPayment(testAuctionID, testWinnerID, 5000)

	if !payment.IsOwnedBy(testWinnerID) {
		t.Error("expected IsOwnedBy to return true for winner")
	}
	if payment.IsOwnedBy(uuid.New().String()) {
		t.Error("expected IsOwnedBy to return false for non-winner")
	}
}

func TestPayment_Refund(t *testing.T) {
	payment, _ := NewPayment(testAuctionID, testWinnerID, 5000)
	payment.ClearEvents()
	_ = payment.Complete()

	if err := payment.Refund("customer request"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if payment.Status() != StatusRefunded {
		t.Errorf("expected status '%s', got '%s'", StatusRefunded, payment.Status())
	}
	if len(payment.Events()) != 2 {
		t.Errorf("expected 2 events (completed + refunded), got %d", len(payment.Events()))
	}
}

func TestPayment_Refund_NotCompleted(t *testing.T) {
	payment, _ := NewPayment(testAuctionID, testWinnerID, 5000)

	if err := payment.Refund("reason"); err == nil {
		t.Error("expected error refunding pending payment")
	}
}

func TestReconstructPayment(t *testing.T) {
	id := uuid.New().String()
	now := time.Now()
	payment := ReconstructPayment(id, testAuctionID, testWinnerID, 5000, StatusCompleted, now, now)

	if payment.ID() != id {
		t.Errorf("expected ID '%s', got '%s'", id, payment.ID())
	}
	if payment.Status() != StatusCompleted {
		t.Errorf("expected status '%s', got '%s'", StatusCompleted, payment.Status())
	}
	if len(payment.Events()) != 0 {
		t.Error("reconstructed payment should have no events")
	}
}
