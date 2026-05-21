package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

var testSellerID = uuid.New().String()

func futureTime() time.Time {
	return time.Now().Add(24 * time.Hour)
}

func TestNewAuction(t *testing.T) {
	auction, err := NewAuction(testSellerID, "Test Auction", "Description", 1000, futureTime())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := uuid.Parse(auction.ID()); err != nil {
		t.Errorf("expected valid UUID for ID, got '%s'", auction.ID())
	}
	if auction.SellerID() != testSellerID {
		t.Errorf("expected SellerID '%s', got '%s'", testSellerID, auction.SellerID())
	}
	if auction.Title() != "Test Auction" {
		t.Errorf("expected title 'Test Auction', got '%s'", auction.Title())
	}
	if auction.Status() != StatusDraft {
		t.Errorf("expected status '%s', got '%s'", StatusDraft, auction.Status())
	}
	if auction.StartPrice() != 1000 {
		t.Errorf("expected start price 1000, got %d", auction.StartPrice())
	}
}

func TestNewAuction_Invariants(t *testing.T) {
	tests := []struct {
		name     string
		sellerID string
		title    string
	}{
		{"empty sellerID", "", "Title"},
		{"empty title", testSellerID, ""},
		{"both empty", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAuction(tt.sellerID, tt.title, "", 1000, futureTime())
			if err == nil {
				t.Errorf("expected error for %s", tt.name)
			}
		})
	}
}

func TestAuction_IsOwnedBy(t *testing.T) {
	auction, _ := NewAuction(testSellerID, "Test", "", 100, futureTime())

	if !auction.IsOwnedBy(testSellerID) {
		t.Error("expected IsOwnedBy to return true for owner")
	}
	if auction.IsOwnedBy(uuid.New().String()) {
		t.Error("expected IsOwnedBy to return false for non-owner")
	}
}

func TestAuction_Open(t *testing.T) {
	auction, _ := NewAuction(testSellerID, "Test", "", 100, futureTime())
	auction.ClearEvents()

	if err := auction.Open(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auction.Status() != StatusOpen {
		t.Errorf("expected status '%s', got '%s'", StatusOpen, auction.Status())
	}
	if len(auction.Events()) != 1 {
		t.Errorf("expected 1 event, got %d", len(auction.Events()))
	}

	if err := auction.Open(); err == nil {
		t.Error("expected error when opening already open auction")
	}
}

func TestAuction_Close(t *testing.T) {
	auction, _ := NewAuction(testSellerID, "Test", "", 100, futureTime())

	if err := auction.Close(); err == nil {
		t.Error("expected error when closing draft auction")
	}

	if err := auction.Open(); err != nil {
		t.Fatal(err)
	}
	auction.ClearEvents()

	if err := auction.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auction.Status() != StatusClosed {
		t.Errorf("expected status '%s', got '%s'", StatusClosed, auction.Status())
	}
	if len(auction.Events()) != 1 {
		t.Errorf("expected 1 event, got %d", len(auction.Events()))
	}
}

func TestAuction_Settle(t *testing.T) {
	auction, _ := NewAuction(testSellerID, "Test", "", 100, futureTime())

	if err := auction.Settle(); err == nil {
		t.Error("expected error when settling draft auction")
	}

	if err := auction.Open(); err != nil {
		t.Fatal(err)
	}
	if err := auction.Close(); err != nil {
		t.Fatal(err)
	}
	auction.ClearEvents()

	if err := auction.Settle(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auction.Status() != StatusSettled {
		t.Errorf("expected status '%s', got '%s'", StatusSettled, auction.Status())
	}
}

func TestAuction_Cancel(t *testing.T) {
	auction, _ := NewAuction(testSellerID, "Test", "", 100, futureTime())

	if err := auction.Cancel(); err != nil {
		t.Fatalf("unexpected error cancelling draft: %v", err)
	}
	if auction.Status() != StatusCancelled {
		t.Errorf("expected status '%s', got '%s'", StatusCancelled, auction.Status())
	}
}

func TestAuction_Cancel_FromClosed(t *testing.T) {
	auction, _ := NewAuction(testSellerID, "Test", "", 100, futureTime())
	if err := auction.Open(); err != nil {
		t.Fatal(err)
	}
	if err := auction.Close(); err != nil {
		t.Fatal(err)
	}
	auction.ClearEvents()

	if err := auction.Cancel(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if auction.Status() != StatusCancelled {
		t.Errorf("expected status '%s', got '%s'", StatusCancelled, auction.Status())
	}
}

func TestAuction_Cancel_FromOpen_Fails(t *testing.T) {
	auction, _ := NewAuction(testSellerID, "Test", "", 100, futureTime())
	if err := auction.Open(); err != nil {
		t.Fatal(err)
	}

	if err := auction.Cancel(); err == nil {
		t.Error("expected error when cancelling open auction")
	}
}

func TestAuction_ClearEvents(t *testing.T) {
	auction, _ := NewAuction(testSellerID, "Test", "", 100, futureTime())
	if err := auction.Open(); err != nil {
		t.Fatal(err)
	}

	if len(auction.Events()) == 0 {
		t.Error("expected events after Open")
	}

	auction.ClearEvents()
	if len(auction.Events()) != 0 {
		t.Error("expected empty events after ClearEvents")
	}
}

func TestReconstructAuction(t *testing.T) {
	id := uuid.New().String()
	now := time.Now()
	endTime := now.Add(24 * time.Hour)

	auction := ReconstructAuction(id, testSellerID, "Title", "Desc", 500, StatusOpen, endTime, now, now)

	if auction.ID() != id {
		t.Errorf("expected ID '%s', got '%s'", id, auction.ID())
	}
	if auction.Status() != StatusOpen {
		t.Errorf("expected status '%s', got '%s'", StatusOpen, auction.Status())
	}
	if len(auction.Events()) != 0 {
		t.Error("reconstructed auction should have no events")
	}
}
