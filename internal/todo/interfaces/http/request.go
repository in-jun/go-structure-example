package http

import (
	"time"

	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

type CreateTodoRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
}

type UpdateTodoRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
}

type UpdateTodoStatusRequest struct {
	Status entity.Status `json:"status"`
}
