package vo

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func futureTime() time.Time {
	return time.Now().Add(24 * time.Hour)
}

func TestNewAuctionIDVO(t *testing.T) {
	valid := uuid.New().String()
	if _, err := NewAuctionIDVO(valid); err != nil {
		t.Errorf("expected no error for valid UUID, got %v", err)
	}

	if _, err := NewAuctionIDVO("invalid"); err == nil {
		t.Error("expected error for invalid UUID")
	}

	if _, err := NewAuctionIDVO(""); err == nil {
		t.Error("expected error for empty string")
	}
}

func TestNewSellerIDVO(t *testing.T) {
	valid := uuid.New().String()
	if _, err := NewSellerIDVO(valid); err != nil {
		t.Errorf("expected no error for valid UUID, got %v", err)
	}

	if _, err := NewSellerIDVO("bad"); err == nil {
		t.Error("expected error for invalid UUID")
	}
}

func TestNewCreateVO(t *testing.T) {
	if _, err := NewCreateVO("Title", "Desc", 100, futureTime()); err != nil {
		t.Errorf("expected no error for valid input, got %v", err)
	}

	if _, err := NewCreateVO("", "Desc", 100, futureTime()); err == nil {
		t.Error("expected error for empty title")
	}

	if _, err := NewCreateVO("Title", "Desc", 0, futureTime()); err == nil {
		t.Error("expected error for zero price")
	}

	if _, err := NewCreateVO("Title", "Desc", -1, futureTime()); err == nil {
		t.Error("expected error for negative price")
	}
}
