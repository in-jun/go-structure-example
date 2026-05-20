package entity

import (
	"testing"
	"time"
)

func TestNewTodo(t *testing.T) {
	tests := []struct {
		name        string
		userID      uint
		title       string
		description string
		wantErr     bool
	}{
		{"valid", 1, "Buy groceries", "Milk and eggs", false},
		{"zero userID", 0, "Buy groceries", "", true},
		{"empty title", 1, "", "description", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := NewTodo(tt.userID, tt.title, tt.description, time.Now().Add(time.Hour))
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if todo.Title() != tt.title {
					t.Errorf("expected title %q, got %q", tt.title, todo.Title())
				}
				if todo.Status() != StatusPending {
					t.Errorf("expected status pending, got %q", todo.Status())
				}
			}
		})
	}
}

func TestTodo_SetStatus(t *testing.T) {
	todo, _ := NewTodo(1, "Test", "", time.Now().Add(time.Hour))
	if todo.Status() != StatusPending {
		t.Errorf("expected pending, got %q", todo.Status())
	}

	todo.SetStatus(StatusCompleted)
	if todo.Status() != StatusCompleted {
		t.Errorf("expected completed, got %q", todo.Status())
	}
}

func TestTodo_Update(t *testing.T) {
	todo, _ := NewTodo(1, "Original", "Desc", time.Now().Add(time.Hour))
	newDue := time.Now().Add(24 * time.Hour)

	todo.Update("Updated", "New desc", newDue)

	if todo.Title() != "Updated" {
		t.Errorf("expected Updated, got %q", todo.Title())
	}
	if todo.Description() != "New desc" {
		t.Errorf("expected New desc, got %q", todo.Description())
	}
}
