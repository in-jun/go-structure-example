package event

import "testing"

const testPaymentID = "test-payment-id"
const testAuctionID = "test-auction-id"
const testWinnerID = "test-winner-id"

func TestPaymentCreated_EventName(t *testing.T) {
	e := NewPaymentCreated(testPaymentID, testAuctionID, testWinnerID, 5000)
	if e.EventName() != "payment.created" {
		t.Errorf("EventName = %q, want payment.created", e.EventName())
	}
	if e.AggregateID() != testPaymentID {
		t.Errorf("AggregateID = %q, want %q", e.AggregateID(), testPaymentID)
	}
	if e.OccurredAt().IsZero() {
		t.Error("OccurredAt should not be zero")
	}
}

func TestPaymentCompleted_EventName(t *testing.T) {
	e := NewPaymentCompleted(testPaymentID, testAuctionID, testWinnerID, 5000)
	if e.EventName() != "payment.completed" {
		t.Errorf("EventName = %q, want payment.completed", e.EventName())
	}
	if e.AggregateID() != testPaymentID {
		t.Errorf("AggregateID = %q, want %q", e.AggregateID(), testPaymentID)
	}
}

func TestPaymentFailed_EventName(t *testing.T) {
	e := NewPaymentFailed(testPaymentID, testAuctionID, testWinnerID, 5000, "declined")
	if e.EventName() != "payment.failed" {
		t.Errorf("EventName = %q, want payment.failed", e.EventName())
	}
	if e.AggregateID() != testPaymentID {
		t.Errorf("AggregateID = %q, want %q", e.AggregateID(), testPaymentID)
	}
}

func TestPaymentRefunded_EventName(t *testing.T) {
	e := NewPaymentRefunded(testPaymentID, testAuctionID, testWinnerID, 5000, "customer request")
	if e.EventName() != "payment.refunded" {
		t.Errorf("EventName = %q, want payment.refunded", e.EventName())
	}
	if e.AggregateID() != testPaymentID {
		t.Errorf("AggregateID = %q, want %q", e.AggregateID(), testPaymentID)
	}
}
