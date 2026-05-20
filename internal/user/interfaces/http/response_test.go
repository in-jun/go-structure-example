package http

import (
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/user/application/query"
)

func TestToUserResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	result := &query.Result{
		ID:        testUUID,
		Email:     "test@example.com",
		Name:      "Test User",
		CreatedAt: now,
	}

	resp := toUserResponse(result)

	if resp.ID != testUUID {
		t.Errorf("ID = %q, want %q", resp.ID, testUUID)
	}
	if resp.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", resp.Email, "test@example.com")
	}
	if resp.Name != "Test User" {
		t.Errorf("Name = %q, want %q", resp.Name, "Test User")
	}
	if !resp.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v, want %v", resp.CreatedAt, now)
	}
}
