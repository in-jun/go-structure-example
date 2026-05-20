package vo

import (
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

func TestNewCreateTodoVO(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{"valid", "Buy groceries", false},
		{"empty title", "", true},
		{"too long title", string(make([]byte, 256)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewCreateTodoVO(tt.title, "desc", time.Now().Add(time.Hour))
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if v.Title != tt.title {
					t.Errorf("expected title %q, got %q", tt.title, v.Title)
				}
			}
		})
	}
}

func TestNewUpdateTodoVO(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{"valid", "Updated title", false},
		{"empty title", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUpdateTodoVO(tt.title, "", time.Now().Add(time.Hour))
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestNewUpdateStatusVO(t *testing.T) {
	tests := []struct {
		name    string
		status  entity.Status
		wantErr bool
	}{
		{"pending", entity.StatusPending, false},
		{"completed", entity.StatusCompleted, false},
		{"invalid", entity.Status("invalid"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUpdateStatusVO(tt.status)
			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
