package entity

import (
	"testing"
	"time"
)

const testUUID = "550e8400-e29b-41d4-a716-446655440000"
const testUUID2 = "660e8400-e29b-41d4-a716-446655440000"

func TestNewTodo(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		title       string
		description string
		wantError   bool
	}{
		{"valid", testUUID, "Buy groceries", "Milk and eggs", false},
		{"empty userID", "", "Buy groceries", "", true},
		{"empty title", testUUID, "", "description", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := NewTodo(tt.userID, tt.title, tt.description, time.Now().Add(time.Hour))
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", todo)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if !tt.wantError && todo.Status() != StatusPending {
				t.Errorf("expected status pending, got %q", todo.Status())
			}
			if !tt.wantError && todo.ID() == "" {
				t.Error("expected non-empty UUID for new todo")
			}
		})
	}
}

func TestTodo_SetStatus(t *testing.T) {
	todo, _ := NewTodo(testUUID, "Test", "", time.Now().Add(time.Hour))
	if todo.Status() != StatusPending {
		t.Errorf("expected pending, got %q", todo.Status())
	}

	todo.SetStatus(StatusCompleted)
	if todo.Status() != StatusCompleted {
		t.Errorf("expected completed, got %q", todo.Status())
	}
}

func TestTodo_Update(t *testing.T) {
	todo, _ := NewTodo(testUUID, "Original", "Desc", time.Now().Add(time.Hour))
	newDue := time.Now().Add(24 * time.Hour)

	todo.Update("Updated", "New desc", newDue)

	if todo.Title() != "Updated" {
		t.Errorf("expected Updated, got %q", todo.Title())
	}
	if todo.Description() != "New desc" {
		t.Errorf("expected New desc, got %q", todo.Description())
	}
}

func TestReconstructTodo(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		id        string
		userID    string
		title     string
		wantError bool
	}{
		{"valid", testUUID, testUUID2, "Test Todo", false},
		{"empty id", "", testUUID2, "Test Todo", true},
		{"empty userID", testUUID, "", "Test Todo", true},
		{"empty title", testUUID, testUUID2, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := ReconstructTodo(tt.id, tt.userID, tt.title, "", StatusPending, now.Add(time.Hour), now, now)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", todo)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
