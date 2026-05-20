package http

import (
	"time"

	"github.com/in-jun/go-structure-example/internal/todo/application/query"
	"github.com/in-jun/go-structure-example/internal/todo/domain/entity"
)

type TodoResponse struct {
	ID          uint          `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      entity.Status `json:"status"`
	DueDate     time.Time     `json:"due_date"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type TodoListResponse struct {
	Todos []TodoResponse `json:"todos"`
	Total int64          `json:"total"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func toTodoResponse(r *query.Result) TodoResponse {
	return TodoResponse{
		ID:          r.ID,
		Title:       r.Title,
		Description: r.Description,
		Status:      r.Status,
		DueDate:     r.DueDate,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func toTodoListResponse(r *query.ListResult) *TodoListResponse {
	todos := make([]TodoResponse, 0, len(r.Todos))
	for i := range r.Todos {
		todos = append(todos, toTodoResponse(&r.Todos[i]))
	}
	return &TodoListResponse{Todos: todos, Total: r.Total}
}
