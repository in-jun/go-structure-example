package vo

import (
	"testing"
	"time"

	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

func TestNewCreateTodoVO(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		wantError bool
	}{
		{"valid", "Buy groceries", false},
		{"empty title", "", true},
		{"too long title", string(make([]byte, 256)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := NewCreateTodoVO(tt.title, "desc", time.Now().Add(time.Hour))
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", vo)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestNewUpdateTodoVO(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		wantError bool
	}{
		{"valid", "Updated title", false},
		{"empty title", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := NewUpdateTodoVO(tt.title, "", time.Now().Add(time.Hour))
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", vo)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}

func TestNewUpdateStatusVO(t *testing.T) {
	tests := []struct {
		name      string
		status    entity.Status
		wantError bool
	}{
		{"pending", entity.StatusPending, false},
		{"completed", entity.StatusCompleted, false},
		{"invalid", entity.Status("invalid"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vo, err := NewUpdateStatusVO(tt.status)
			if tt.wantError && err == nil {
				t.Errorf("expected error, got %+v", vo)
			}
			if !tt.wantError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
		})
	}
}
